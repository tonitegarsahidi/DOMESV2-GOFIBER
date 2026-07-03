package repository

import (
	"strings"

	"domesv2/config/database"
	"domesv2/internal/model"
	"domesv2/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DocumentRepository interface {
	Create(doc *model.Document) error
	GetByID(id string) (*model.Document, error)
	GetBySlug(slug string) (*model.Document, error)
	Update(doc *model.Document) error
	Delete(id string) error
	ListPublic(filters map[string]interface{}) ([]model.Document, int, error)
	ListSubmissions(status string, search string, page int, limit int) ([]model.Document, int, error)
	GetRelated(doc *model.Document) ([]model.Document, error)
	GetGlobalStats() (map[string]interface{}, error)
	GetOverviewAnalytics() (map[string]interface{}, error)
	GetUploadsOverTime(fromYear, toYear int) ([]map[string]interface{}, error)
	GetBySdgAnalytics() ([]map[string]interface{}, error)
	GetByAgencyAnalytics() ([]map[string]interface{}, error)
	GetBySectorAnalytics() ([]map[string]interface{}, error)
	GetByLanguageAnalytics() ([]map[string]interface{}, error)
}

type documentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository() DocumentRepository {
	return &documentRepository{
		db: database.GetDB(),
	}
}

func (r *documentRepository) Create(doc *model.Document) error {
	if err := r.db.Create(doc).Error; err != nil {
		zap.L().Error("Failed to create document", zap.Error(err))
		return errors.NewInternalServerError("Failed to create submission", "DATABASE_ERROR")
	}
	return nil
}

func (r *documentRepository) GetByID(id string) (*model.Document, error) {
	var doc model.Document
	err := r.db.Preload("LeadAgency").Preload("JointProgramme").
		Preload("Sdgs").Preload("Sectors").Preload("Lnobs").
		Preload("Author").First(&doc, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("Document not found", "DOCUMENT_NOT_FOUND")
		}
		zap.L().Error("Failed to fetch document by id", zap.Error(err))
		return nil, errors.NewInternalServerError("Failed to fetch document", "DATABASE_ERROR")
	}
	return &doc, nil
}

func (r *documentRepository) GetBySlug(slug string) (*model.Document, error) {
	var doc model.Document
	err := r.db.Preload("LeadAgency").Preload("JointProgramme").
		Preload("Sdgs").Preload("Sectors").Preload("Lnobs").
		Preload("Author").Where("slug = ?", slug).First(&doc).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("Document not found", "DOCUMENT_NOT_FOUND")
		}
		zap.L().Error("Failed to fetch document by slug", zap.Error(err))
		return nil, errors.NewInternalServerError("Failed to fetch document", "DATABASE_ERROR")
	}
	return &doc, nil
}

func (r *documentRepository) Update(doc *model.Document) error {
	if err := r.db.Save(doc).Error; err != nil {
		zap.L().Error("Failed to update document", zap.Error(err))
		return errors.NewInternalServerError("Failed to update document", "DATABASE_ERROR")
	}

	// Update many-to-many associations explicitly to prevent orphans/duplicates
	r.db.Model(doc).Association("Sdgs").Replace(doc.Sdgs)
	r.db.Model(doc).Association("Sectors").Replace(doc.Sectors)
	r.db.Model(doc).Association("Lnobs").Replace(doc.Lnobs)

	return nil
}

func (r *documentRepository) Delete(id string) error {
	var doc model.Document
	if err := r.db.First(&doc, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError("Submission not found", "SUBMISSION_NOT_FOUND")
		}
		return errors.NewInternalServerError("Database query error", "DATABASE_ERROR")
	}

	// Clear associations first
	r.db.Model(&doc).Association("Sdgs").Clear()
	r.db.Model(&doc).Association("Sectors").Clear()
	r.db.Model(&doc).Association("Lnobs").Clear()

	if err := r.db.Delete(&doc).Error; err != nil {
		zap.L().Error("Failed to delete document", zap.Error(err))
		return errors.NewInternalServerError("Failed to delete submission", "DATABASE_ERROR")
	}
	return nil
}

func (r *documentRepository) ListPublic(filters map[string]interface{}) ([]model.Document, int, error) {
	var docs []model.Document
	query := r.db.Model(&model.Document{}).Where("V2Documents.status = ? AND V2Documents.isActive = ?", "published", true)

	// Pre-resolve language terms if language filter is active to handle comma-separated values and code/name mapping
	var langConditions []string
	var langArgs []interface{}
	if langsStr, ok := filters["langs"].(string); ok && langsStr != "" {
		langs := strings.Split(langsStr, ",")
		var dbLangs []model.Language
		if err := r.db.Where("code IN (?) OR LOWER(name) IN (?)", langs, langs).Find(&dbLangs).Error; err == nil {
			searchTerms := make(map[string]bool)
			for _, lang := range langs {
				searchTerms[strings.ToLower(lang)] = true
			}
			for _, l := range dbLangs {
				searchTerms[strings.ToLower(l.Code)] = true
				searchTerms[strings.ToLower(l.Name)] = true
			}
			for term := range searchTerms {
				langConditions = append(langConditions, "LOWER(V2Documents.language) LIKE ?")
				langArgs = append(langArgs, "%"+term+"%")
			}
		}
	}

	// Apply filter: Search Text
	if q, ok := filters["q"].(string); ok && q != "" {
		searchText := "%" + strings.ToLower(q) + "%"
		query = query.Where("LOWER(V2Documents.title) LIKE ? OR LOWER(V2Documents.description) LIKE ? OR LOWER(V2Documents.abstract) LIKE ?", searchText, searchText, searchText)
	}

	// Apply filter: Agencies (comma separated)
	if agenciesStr, ok := filters["agencies"].(string); ok && agenciesStr != "" {
		agencies := strings.Split(agenciesStr, ",")
		query = query.Where("V2Documents.lead_agency_code IN (?)", agencies)
	}

	// Apply filter: Joint programmes
	if jpStr, ok := filters["jointProgrammes"].(string); ok && jpStr != "" {
		jps := strings.Split(jpStr, ",")
		query = query.Where("V2Documents.joint_programme_code IN (?)", jps)
	}

	// Apply filter: Years
	if yearFrom, ok := filters["yearFrom"].(int); ok && yearFrom > 0 {
		query = query.Where("V2Documents.year >= ?", yearFrom)
	}
	if yearTo, ok := filters["yearTo"].(int); ok && yearTo > 0 {
		query = query.Where("V2Documents.year <= ?", yearTo)
	}

	// Apply filter: Languages
	if len(langConditions) > 0 {
		query = query.Where(strings.Join(langConditions, " OR "), langArgs...)
	}

	// Apply filter: SDGs (many-to-many relation filter)
	if sdgsStr, ok := filters["sdgs"].(string); ok && sdgsStr != "" {
		sdgs := strings.Split(sdgsStr, ",")
		query = query.Joins("JOIN v2_document_sdgs ON v2_document_sdgs.document_id = V2Documents.id").
			Joins("JOIN V2MasterSdgs ON V2MasterSdgs.id = v2_document_sdgs.sdg_id").
			Where("V2MasterSdgs.code IN (?)", sdgs).
			Group("V2Documents.id")
	}

	// Apply filter: Sectors (many-to-many relation filter)
	if sectorsStr, ok := filters["sectors"].(string); ok && sectorsStr != "" {
		sectors := strings.Split(sectorsStr, ",")
		query = query.Joins("JOIN v2_document_sectors ON v2_document_sectors.document_id = V2Documents.id").
			Joins("JOIN V2MasterSectors ON V2MasterSectors.id = v2_document_sectors.sector_id").
			Where("V2MasterSectors.code IN (?)", sectors).
			Group("V2Documents.id")
	}

	// Apply filter: LNOBs (many-to-many relation filter)
	if lnobsStr, ok := filters["lnobs"].(string); ok && lnobsStr != "" {
		lnobs := strings.Split(lnobsStr, ",")
		query = query.Joins("JOIN v2_document_lnobs ON v2_document_lnobs.document_id = V2Documents.id").
			Joins("JOIN V2MasterLnobs ON V2MasterLnobs.id = v2_document_lnobs.lnob_id").
			Where("V2MasterLnobs.code IN (?)", lnobs).
			Group("V2Documents.id")
	}

	// Apply filter: Non-UN Partners wildcard
	if partnerStr, ok := filters["nonUnPartners"].(string); ok && partnerStr != "" {
		partners := strings.Split(partnerStr, ",")
		var partnerConditions []string
		var partnerArgs []interface{}
		for _, partner := range partners {
			partnerConditions = append(partnerConditions, "LOWER(V2Documents.non_un_partners) LIKE ?")
			partnerArgs = append(partnerArgs, "%"+strings.ToLower(partner)+"%")
		}
		query = query.Where(strings.Join(partnerConditions, " OR "), partnerArgs...)
	}

	// Get total items count (for pagination)
	var totalItems int64
	countQuery := r.db.Model(&model.Document{}).Where("V2Documents.status = ?", "published")
	if q, ok := filters["q"].(string); ok && q != "" {
		searchText := "%" + strings.ToLower(q) + "%"
		countQuery = countQuery.Where("LOWER(V2Documents.title) LIKE ? OR LOWER(V2Documents.description) LIKE ? OR LOWER(V2Documents.abstract) LIKE ?", searchText, searchText, searchText)
	}
	if agenciesStr, ok := filters["agencies"].(string); ok && agenciesStr != "" {
		countQuery = countQuery.Where("V2Documents.lead_agency_code IN (?)", strings.Split(agenciesStr, ","))
	}
	if jpStr, ok := filters["jointProgrammes"].(string); ok && jpStr != "" {
		countQuery = countQuery.Where("V2Documents.joint_programme_code IN (?)", strings.Split(jpStr, ","))
	}
	if yearFrom, ok := filters["yearFrom"].(int); ok && yearFrom > 0 {
		countQuery = countQuery.Where("V2Documents.year >= ?", yearFrom)
	}
	if yearTo, ok := filters["yearTo"].(int); ok && yearTo > 0 {
		countQuery = countQuery.Where("V2Documents.year <= ?", yearTo)
	}
	if len(langConditions) > 0 {
		countQuery = countQuery.Where(strings.Join(langConditions, " OR "), langArgs...)
	}
	if sdgsStr, ok := filters["sdgs"].(string); ok && sdgsStr != "" {
		countQuery = countQuery.Joins("JOIN v2_document_sdgs ON v2_document_sdgs.document_id = V2Documents.id").
			Joins("JOIN V2MasterSdgs ON V2MasterSdgs.id = v2_document_sdgs.sdg_id").
			Where("V2MasterSdgs.code IN (?)", strings.Split(sdgsStr, ","))
	}
	if sectorsStr, ok := filters["sectors"].(string); ok && sectorsStr != "" {
		countQuery = countQuery.Joins("JOIN v2_document_sectors ON v2_document_sectors.document_id = V2Documents.id").
			Joins("JOIN V2MasterSectors ON V2MasterSectors.id = v2_document_sectors.sector_id").
			Where("V2MasterSectors.code IN (?)", strings.Split(sectorsStr, ","))
	}
	if lnobsStr, ok := filters["lnobs"].(string); ok && lnobsStr != "" {
		countQuery = countQuery.Joins("JOIN v2_document_lnobs ON v2_document_lnobs.document_id = V2Documents.id").
			Joins("JOIN V2MasterLnobs ON V2MasterLnobs.id = v2_document_lnobs.lnob_id").
			Where("V2MasterLnobs.code IN (?)", strings.Split(lnobsStr, ","))
	}

	countQuery.Select("COUNT(DISTINCT V2Documents.id)").Count(&totalItems)

	// Apply Sorting
	sort := "newest"
	if s, ok := filters["sort"].(string); ok && s != "" {
		sort = s
	}
	switch sort {
	case "oldest":
		query = query.Order("V2Documents.year asc, V2Documents.createdAt asc")
	case "popular":
		query = query.Order("V2Documents.views desc, V2Documents.downloads desc")
	case "relevance":
		query = query.Order("V2Documents.createdAt desc")
	case "newest":
		fallthrough
	default:
		query = query.Order("V2Documents.year desc, V2Documents.createdAt desc")
	}

	// Apply Pagination
	page := 1
	limit := 12
	if p, ok := filters["page"].(int); ok && p > 0 {
		page = p
	}
	if l, ok := filters["limit"].(int); ok && l > 0 {
		limit = l
	}
	offset := (page - 1) * limit
	query = query.Limit(limit).Offset(offset)

	// Fetch documents with relations
	err := query.Preload("LeadAgency").Preload("JointProgramme").
		Preload("Sdgs").Preload("Sectors").Preload("Lnobs").
		Preload("Author").Find(&docs).Error
	if err != nil {
		zap.L().Error("Failed to list public documents", zap.Error(err))
		return nil, 0, errors.NewInternalServerError("Failed to retrieve documents", "DATABASE_ERROR")
	}

	return docs, int(totalItems), nil
}

func (r *documentRepository) ListSubmissions(status string, search string, page int, limit int) ([]model.Document, int, error) {
	var docs []model.Document
	query := r.db.Model(&model.Document{})

	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	if search != "" {
		query = query.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(search)+"%")
	}

	var totalItems int64
	query.Count(&totalItems)

	offset := (page - 1) * limit
	err := query.Order("createdAt desc").Limit(limit).Offset(offset).
		Preload("LeadAgency").Preload("Author").Find(&docs).Error
	if err != nil {
		zap.L().Error("Failed to list submissions for CMS", zap.Error(err))
		return nil, 0, errors.NewInternalServerError("Failed to retrieve submissions", "DATABASE_ERROR")
	}

	return docs, int(totalItems), nil
}

func (r *documentRepository) GetRelated(doc *model.Document) ([]model.Document, error) {
	var relatedDocs []model.Document
	var sdgCodes []string
	for _, sdg := range doc.Sdgs {
		sdgCodes = append(sdgCodes, sdg.Code)
	}

	if len(sdgCodes) == 0 {
		// Fallback to same agency
		err := r.db.Where("lead_agency_code = ? AND id != ? AND status = ?", doc.LeadAgencyCode, doc.ID, "published").
			Limit(3).Preload("Sdgs").Find(&relatedDocs).Error
		return relatedDocs, err
	}

	err := r.db.Model(&model.Document{}).
		Joins("JOIN v2_document_sdgs ON v2_document_sdgs.document_id = V2Documents.id").
		Joins("JOIN V2MasterSdgs ON V2MasterSdgs.id = v2_document_sdgs.sdg_id").
		Where("V2MasterSdgs.code IN (?) AND V2Documents.id != ? AND V2Documents.status = ?", sdgCodes, doc.ID, "published").
		Group("V2Documents.id").
		Limit(3).
		Preload("Sdgs").
		Find(&relatedDocs).Error
	if err != nil {
		zap.L().Error("Failed to fetch related documents", zap.Error(err))
		return nil, errors.NewInternalServerError("Failed to fetch related documents", "DATABASE_ERROR")
	}
	return relatedDocs, nil
}

func (r *documentRepository) GetGlobalStats() (map[string]interface{}, error) {
	var totalDocs int64
	var totalAgencies int64
	var totalPartners int64 = 35
	var totalSdgGoals int64 = 17

	r.db.Model(&model.Document{}).Where("status = ?", "published").Count(&totalDocs)
	r.db.Model(&model.Agency{}).Count(&totalAgencies)

	return map[string]interface{}{
		"total_documents": int(totalDocs),
		"total_agencies":  int(totalAgencies),
		"total_partners":  int(totalPartners),
		"total_sdg_goals": int(totalSdgGoals),
	}, nil
}

func (r *documentRepository) GetOverviewAnalytics() (map[string]interface{}, error) {
	var totalDocs int64
	var activeAgencies int64
	var monthlyDownloads int64 = 84200
	var totalViews int64 = 456000
	var totalDownloads int64 = 189000

	r.db.Model(&model.Document{}).Where("status = ?", "published").Count(&totalDocs)
	r.db.Model(&model.Agency{}).Count(&activeAgencies)

	var stats struct {
		Views     int64
		Downloads int64
	}
	r.db.Model(&model.Document{}).Select("SUM(views) as views, SUM(downloads) as downloads").Scan(&stats)
	if stats.Views > 0 {
		totalViews = stats.Views
	}
	if stats.Downloads > 0 {
		totalDownloads = stats.Downloads
	}

	return map[string]interface{}{
		"total_documents":   int(totalDocs),
		"active_agencies":   int(activeAgencies),
		"monthly_downloads": int(monthlyDownloads),
		"total_views":       int(totalViews),
		"total_downloads":   int(totalDownloads),
	}, nil
}

func (r *documentRepository) GetUploadsOverTime(fromYear, toYear int) ([]map[string]interface{}, error) {
	var results []struct {
		Year  int
		Count int64
	}

	err := r.db.Model(&model.Document{}).
		Select("year, count(id) as count").
		Where("year >= ? AND year <= ? AND status = ?", fromYear, toYear, "published").
		Group("year").
		Order("year asc").
		Scan(&results).Error

	if err != nil {
		zap.L().Error("Failed to fetch uploads over time analytics", zap.Error(err))
		return nil, errors.NewInternalServerError("Failed to calculate analytics", "DATABASE_ERROR")
	}

	var list []map[string]interface{}
	yearMap := make(map[int]int64)
	for _, res := range results {
		yearMap[res.Year] = res.Count
	}

	for y := fromYear; y <= toYear; y++ {
		cnt := yearMap[y]
		list = append(list, map[string]interface{}{
			"year":  y,
			"count": cnt,
		})
	}
	return list, nil
}

func (r *documentRepository) GetBySdgAnalytics() ([]map[string]interface{}, error) {
	var results []struct {
		SdgCode string `gorm:"column:sdg_code"`
		Count   int64  `gorm:"column:count"`
	}

	err := r.db.Table("v2_document_sdgs").
		Select("V2MasterSdgs.code as sdg_code, count(v2_document_sdgs.document_id) as count").
		Joins("JOIN V2MasterSdgs ON V2MasterSdgs.id = v2_document_sdgs.sdg_id").
		Group("V2MasterSdgs.code").
		Scan(&results).Error

	if err != nil {
		zap.L().Error("Failed to fetch SDG analytics", zap.Error(err))
		return nil, errors.NewInternalServerError("Failed to calculate analytics", "DATABASE_ERROR")
	}

	var sdgs []model.Sdg
	r.db.Find(&sdgs)
	sdgMap := make(map[string]model.Sdg)
	for _, s := range sdgs {
		sdgMap[s.Code] = s
	}

	var list []map[string]interface{}
	for _, res := range results {
		sdg, exists := sdgMap[res.SdgCode]
		if exists {
			list = append(list, map[string]interface{}{
				"sdg":   sdg.Code,
				"name":  sdg.Name,
				"count": res.Count,
				"color": sdg.Color,
			})
		}
	}

	return list, nil
}

func (r *documentRepository) GetByAgencyAnalytics() ([]map[string]interface{}, error) {
	var results []struct {
		LeadAgencyCode string `gorm:"column:lead_agency_code"`
		Count          int64  `gorm:"column:count"`
	}

	err := r.db.Model(&model.Document{}).
		Select("lead_agency_code, count(id) as count").
		Where("status = ?", "published").
		Group("lead_agency_code").
		Scan(&results).Error

	if err != nil {
		zap.L().Error("Failed to fetch agency analytics", zap.Error(err))
		return nil, errors.NewInternalServerError("Failed to calculate analytics", "DATABASE_ERROR")
	}

	var list []map[string]interface{}
	for _, res := range results {
		if res.LeadAgencyCode != "" {
			list = append(list, map[string]interface{}{
				"agency": res.LeadAgencyCode,
				"count":  res.Count,
			})
		}
	}
	return list, nil
}

func (r *documentRepository) GetBySectorAnalytics() ([]map[string]interface{}, error) {
	var results []struct {
		SectorCode string `gorm:"column:sector_code"`
		Count      int64  `gorm:"column:count"`
	}

	err := r.db.Table("v2_document_sectors").
		Select("V2MasterSectors.code as sector_code, count(v2_document_sectors.document_id) as count").
		Joins("JOIN V2MasterSectors ON V2MasterSectors.id = v2_document_sectors.sector_id").
		Group("V2MasterSectors.code").
		Scan(&results).Error

	if err != nil {
		zap.L().Error("Failed to fetch sector analytics", zap.Error(err))
		return nil, errors.NewInternalServerError("Failed to calculate analytics", "DATABASE_ERROR")
	}

	var sectors []model.Sector
	r.db.Find(&sectors)
	sectorMap := make(map[string]model.Sector)
	for _, s := range sectors {
		sectorMap[s.Code] = s
	}

	var list []map[string]interface{}
	for _, res := range results {
		sec, exists := sectorMap[res.SectorCode]
		if exists {
			list = append(list, map[string]interface{}{
				"sector": sec.Name,
				"count":  res.Count,
			})
		}
	}
	return list, nil
}

func (r *documentRepository) GetByLanguageAnalytics() ([]map[string]interface{}, error) {
	var results []struct {
		Language string `gorm:"column:language"`
		Count    int64  `gorm:"column:count"`
	}

	err := r.db.Model(&model.Document{}).
		Select("language, count(id) as count").
		Where("status = ?", "published").
		Group("language").
		Scan(&results).Error

	if err != nil {
		zap.L().Error("Failed to fetch language analytics", zap.Error(err))
		return nil, errors.NewInternalServerError("Failed to calculate analytics", "DATABASE_ERROR")
	}

	var list []map[string]interface{}
	for _, res := range results {
		if res.Language != "" {
			list = append(list, map[string]interface{}{
				"language": res.Language,
				"count":    res.Count,
			})
		}
	}
	return list, nil
}
