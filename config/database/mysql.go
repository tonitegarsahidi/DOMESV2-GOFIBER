package database

import (
	"fmt"
	"time"

	"domesv2/config"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

const (
	maxRetries       = 5
	initialBackoff   = 1 * time.Second
	backoffMultiplier = 2
)

func InitMySQL(cfg *config.Config) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
		cfg.DB.Charset,
		cfg.DB.ParseTime,
		cfg.DB.Loc,
	)

	var logLevel logger.LogLevel
	if cfg.Server.Env == "development" {
		logLevel = logger.Info
	} else {
		logLevel = logger.Error
	}

	gormConfig := &gorm.Config{
		Logger:                                   logger.Default.LogMode(logLevel),
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	backoff := initialBackoff

	for attempt := 1; attempt <= maxRetries; attempt++ {
		db, err := gorm.Open(mysql.Open(dsn), gormConfig)
		if err == nil {
			// Verify the connection is actually alive
			sqlDB, pingErr := db.DB()
			if pingErr == nil {
				pingErr = sqlDB.Ping()
			}

			if pingErr == nil {
				zap.L().Info("Database connection established successfully",
					zap.Int("attempt", attempt),
				)
				DB = db
				return
			}
			err = pingErr
		}

		if attempt < maxRetries {
			zap.L().Warn("Failed to connect to database, retrying...",
				zap.Int("attempt", attempt),
				zap.Int("max_retries", maxRetries),
				zap.Duration("next_retry_in", backoff),
				zap.Error(err),
			)
			time.Sleep(backoff)
			backoff *= backoffMultiplier
		} else {
			zap.L().Fatal("Failed to connect to database after all retries. Exiting.",
				zap.Int("attempts", maxRetries),
				zap.Error(err),
			)
		}
	}
}

func GetDB() *gorm.DB {
	return DB
}
