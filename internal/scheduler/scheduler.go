package scheduler

import (
	"context"
	"time"

	"domesv2/config/database"
	"domesv2/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		// Run once immediately on start
		zap.L().Info("Running initial document stats synchronization...")
		if err := SyncDocumentStats(); err != nil {
			zap.L().Error("Failed to sync document stats", zap.Error(err))
		}

		for {
			select {
			case <-ticker.C:
				zap.L().Info("Running scheduled document stats synchronization...")
				if err := SyncDocumentStats(); err != nil {
					zap.L().Error("Failed to sync document stats", zap.Error(err))
				}
			case <-ctx.Done():
				ticker.Stop()
				zap.L().Info("Document stats scheduler stopped.")
				return
			}
		}
	}()
}

func SyncDocumentStats() error {
	db := database.GetDB()
	if db == nil {
		return gorm.ErrInvalidDB
	}

	// 1. Get the global checkpoint (max last_processed_log_id across all stats)
	var lastProcessedLogID uint64
	err := db.Model(&model.DocumentStats{}).Select("COALESCE(MAX(last_processed_log_id), 0)").Row().Scan(&lastProcessedLogID)
	if err != nil {
		// If table is empty or doesn't have records, fallback to 0
		lastProcessedLogID = 0
	}

	// 2. Fetch new activities aggregated by document_id and action
	type LogAggregate struct {
		DocumentID string `gorm:"column:document_id"`
		Action     string `gorm:"column:action"`
		Count      int    `gorm:"column:count"`
		MaxLogID   uint64 `gorm:"column:max_log_id"`
	}

	var aggregates []LogAggregate
	err = db.Model(&model.DocumentActivityLog{}).
		Select("document_id, action, COUNT(*) as count, MAX(id) as max_log_id").
		Where("id > ?", lastProcessedLogID).
		Group("document_id, action").
		Scan(&aggregates).Error
	if err != nil {
		return err
	}

	if len(aggregates) == 0 {
		zap.L().Info("No new document activities to sync.")
		return nil
	}

	// 3. Process aggregates inside a database transaction to ensure consistency
	return db.Transaction(func(tx *gorm.DB) error {
		for _, agg := range aggregates {
			// Find the document to retrieve its internal integer ID
			var doc model.Document
			if err := tx.Select("id").Where("uuid = ?", agg.DocumentID).First(&doc).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					// Document was deleted or doesn't exist, skip this log aggregate
					continue
				}
				return err
			}

			// Find or create stats row for this document
			var stats model.DocumentStats
			err := tx.Where("document_id = ?", doc.ID).First(&stats).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					// Create new stats record
					stats = model.DocumentStats{
						DocumentID:         doc.ID,
						TotalViews:         0,
						TotalDownloads:     0,
						LastProcessedLogID: agg.MaxLogID,
					}
					if agg.Action == "view" {
						stats.TotalViews = agg.Count
					} else if agg.Action == "download" {
						stats.TotalDownloads = agg.Count
					}
					if err := tx.Create(&stats).Error; err != nil {
						return err
					}
					continue
				}
				return err
			}

			// Update existing stats record
			if agg.Action == "view" {
				stats.TotalViews += agg.Count
			} else if agg.Action == "download" {
				stats.TotalDownloads += agg.Count
			}

			// Update checkpoint to the maximum log ID processed for this document's batch
			if agg.MaxLogID > stats.LastProcessedLogID {
				stats.LastProcessedLogID = agg.MaxLogID
			}

			if err := tx.Save(&stats).Error; err != nil {
				return err
			}
		}

		zap.L().Info("Successfully synchronized document stats", zap.Int("items_processed", len(aggregates)))
		return nil
	})
}
