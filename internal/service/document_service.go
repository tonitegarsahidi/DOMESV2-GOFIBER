package service

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"domesv2/internal/model"
	"domesv2/internal/repository"
	"domesv2/pkg/errors"
)

type DocumentService interface {
	CreateSubmission(userID uint, req *model.SubmissionRequest) (*model.Document, error)
	SaveDraft(userID uint, submissionID string, step int, data interface{}) (*model.Document, error)
	GetDocumentByID(id string) (*model.DocumentResponse, error)
	GetDocumentBySlug(slug string) (*model.DocumentResponse, error)
	ListPublicDocuments(filters map[string]interface{}) (*model.DocumentListResponse, error)
	SearchPublicDocuments(q string, page int, limit int, sort string, filters map[string]interface{}) (map[string]interface{}, error)
	DeleteSubmission(id string) error
	PublishDocument(id string) (*model.Document, error)
	UnpublishDocument(id string) (*model.Document, error)
	GetRelatedDocuments(id string) ([]model.DocumentResponse, error)
	GenerateDownloadLink(id string) (map[string]interface{}, error)
	GetPlatformStats() (map[string]interface{}, error)
	GetAnalyticsOverview() (map[string]interface{}, error)
	GetUploadsOverTime(fromYear, toYear int) ([]map[string]interface{}, error)
	GetBySdgAnalytics() ([]map[string]interface{}, error)
	GetByAgencyAnalytics() ([]map[string]interface{}, error)
	GetBySectorAnalytics() ([]map[string]interface{}, error)
	GetByLanguageAnalytics() ([]map[string]interface{}, error)
	ListSubmissions(status string, search string, page int, limit int) (map[string]interface{}, error)
	RecordDocumentActivity(docID string, action string, ip string, userAgent string) error
}

type documentService struct {
	docRepo  repository.DocumentRepository
	userRepo repository.UserRepository
}

func NewDocumentService(docRepo repository.DocumentRepository, userRepo repository.UserRepository) DocumentService {
	return &documentService{
		docRepo:  docRepo,
		userRepo: userRepo,
	}
}

func (s *documentService) CreateSubmission(userID uint, req *model.SubmissionRequest) (*model.Document, error) {
	if req.Title == "" {
		return nil, errors.NewValidationError("Title is required", "VALIDATION_FAILED")
	}

	slug := slugify(req.Title)
	// Check unique slug, append timestamp if conflict
	existing, _ := s.docRepo.GetBySlug(slug)
	if existing != nil {
		slug = fmt.Sprintf("%s-%d", slug, time.Now().Unix())
	}

	// Derive year from DateOfPublication
	year := time.Now().Year()
	if req.DateOfPublication != "" {
		parts := strings.Split(req.DateOfPublication, "-")
		if len(parts) > 0 {
			if y, err := strconv.Atoi(parts[0]); err == nil {
				year = y
			}
		}
	}

	// Serialize lists to JSON string
	tagsBytes, _ := json.Marshal(req.Tags)
	thematicBytes, _ := json.Marshal(req.ThematicAreas)
	otherAgenciesBytes, _ := json.Marshal(req.OtherAgencies)
	partnersBytes, _ := json.Marshal(req.NonUnPartners)
	supportingBytes, _ := json.Marshal(req.SupportingFiles)

	// Fetch related entities (Sdgs, Sectors, Lnobs)
	// Let's create database-loaded entities
	var sdgs []model.Sdg
	for _, code := range req.Sdgs {
		sdgs = append(sdgs, model.Sdg{Code: code})
	}
	var sectors []model.Sector
	for _, code := range req.Sectors {
		sectors = append(sectors, model.Sector{Code: code})
	}
	var lnobs []model.Lnob
	for _, code := range req.LnobGroups {
		lnobs = append(lnobs, model.Lnob{Code: code})
	}

	// Generate Code, e.g. UNDP-2024-001 or using agency + year + unique slug prefix
	prefix := "DOC"
	if req.Agency != "" {
		prefix = req.Agency
	}
	code := fmt.Sprintf("%s-%d-%s", prefix, year, strings.ToUpper(generateRandomHex(3)))

	doc := &model.Document{
		Code:                 code,
		Slug:                 slug,
		Title:                req.Title,
		Description:          req.ShortDescription,
		Abstract:             req.Abstract,
		Summary:              req.DetailedSummary,
		Year:                 year,
		DateOfPublication:    req.DateOfPublication,
		TotalPages:           req.TotalPages,
		Language:             req.Language,
		FileURL:              req.FileURL,
		FileSize:             req.FileSize,
		CoverImage:           req.CoverImageURL,
		ExternalURL:          req.ExternalURL,
		Status:               "pending_review", // Status becomes pending_review after step 4 final submit
		FocalPointName:       req.FocalPoint.Name,
		FocalPointEmail:      req.FocalPoint.Email,
		FocalPointPhone:      req.FocalPoint.Phone,
		FocalPointDepartment: req.FocalPoint.Department,
		LeadAgencyCode:       req.Agency,
		JointProgrammeCode:   req.JointProgramme,
		GeographicScope:      req.GeographicScope,
		ThematicAreas:        string(thematicBytes),
		Tags:                 string(tagsBytes),
		OtherAgencies:        string(otherAgenciesBytes),
		NonUnPartners:        string(partnersBytes),
		SupportingFiles:      string(supportingBytes),
		AuthorID:             userID,
		Sdgs:                 sdgs,
		Sectors:              sectors,
		Lnobs:                lnobs,
	}
	if req.IsActive != nil {
		doc.IsActive = req.IsActive
	}

	if err := s.docRepo.Create(doc); err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *documentService) SaveDraft(userID uint, submissionID string, step int, data interface{}) (*model.Document, error) {
	// Let's create or fetch draft submission
	var doc *model.Document
	var err error

	if submissionID != "" {
		doc, err = s.docRepo.GetByID(submissionID)
		if err != nil {
			return nil, err
		}
	} else {
		// Create empty draft
		slug := fmt.Sprintf("draft-%d", time.Now().UnixNano())
		doc = &model.Document{
			Title:    "Draft Submission",
			Slug:     slug,
			Status:   "draft",
			AuthorID: userID,
		}
		if err := s.docRepo.Create(doc); err != nil {
			return nil, err
		}
	}

	// Update draft fields depending on wizard step data
	dataBytes, _ := json.Marshal(data)
	switch step {
	case 2:
		var step2Data struct {
			Title             string        `json:"title"`
			DateOfPublication string        `json:"date_of_publication"`
			TotalPages        int           `json:"total_pages"`
			Language          string        `json:"language"`
			PublicationStatus string        `json:"publication_status"`
			ShortSummary      string        `json:"short_summary"`
			Tags              []string      `json:"tags"`
			FocalPointName    string        `json:"focal_point_name"`
			FocalPointEmail   string        `json:"focal_point_email"`
			FocalPointPhone   string        `json:"focal_point_phone"`
			FocalPointDept    string        `json:"focal_point_department"`
			IsActive          *bool         `json:"is_active"`
		}
		if err := json.Unmarshal(dataBytes, &step2Data); err == nil {
			if step2Data.Title != "" {
				doc.Title = step2Data.Title
				doc.Slug = slugify(step2Data.Title)
			}
			doc.DateOfPublication = step2Data.DateOfPublication
			doc.TotalPages = step2Data.TotalPages
			doc.Language = step2Data.Language
			doc.Description = step2Data.ShortSummary
			doc.FocalPointName = step2Data.FocalPointName
			doc.FocalPointEmail = step2Data.FocalPointEmail
			doc.FocalPointPhone = step2Data.FocalPointPhone
			doc.FocalPointDepartment = step2Data.FocalPointDept
			tagsBytes, _ := json.Marshal(step2Data.Tags)
			doc.Tags = string(tagsBytes)
			if step2Data.IsActive != nil {
				doc.IsActive = step2Data.IsActive
			}
		}
	case 3:
		var step3Data struct {
			Abstract        string                 `json:"abstract"`
			DetailedSummary string                 `json:"detailed_summary"`
			FileURL         string                 `json:"file_url"`
			FileSize        string                 `json:"file_size"`
			CoverImageURL   string                 `json:"cover_image_url"`
			SupportingFiles []model.SupportingFile `json:"supporting_files"`
		}
		if err := json.Unmarshal(dataBytes, &step3Data); err == nil {
			doc.Abstract = step3Data.Abstract
			doc.Summary = step3Data.DetailedSummary
			doc.FileURL = step3Data.FileURL
			doc.FileSize = step3Data.FileSize
			doc.CoverImage = step3Data.CoverImageURL
			supBytes, _ := json.Marshal(step3Data.SupportingFiles)
			doc.SupportingFiles = string(supBytes)
		}
	}

	doc.Status = "draft"
	if err := s.docRepo.Update(doc); err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *documentService) GetDocumentByID(id string) (*model.DocumentResponse, error) {
	doc, err := s.docRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Increment view count asynchronously
	go func(d *model.Document) {
		d.Views += 1
		s.docRepo.Update(d)
	}(doc)

	resp := mapToDocumentResponse(doc)
	return &resp, nil
}

func (s *documentService) GetDocumentBySlug(slug string) (*model.DocumentResponse, error) {
	doc, err := s.docRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	go func(d *model.Document) {
		d.Views += 1
		s.docRepo.Update(d)
	}(doc)

	resp := mapToDocumentResponse(doc)
	return &resp, nil
}

func (s *documentService) ListPublicDocuments(filters map[string]interface{}) (*model.DocumentListResponse, error) {
	docs, totalItems, err := s.docRepo.ListPublic(filters)
	if err != nil {
		return nil, err
	}

	var items []model.DocumentListItem
	for _, doc := range docs {
		var sdgs []model.SdgDTO
		for _, s := range doc.Sdgs {
			sdgs = append(sdgs, model.SdgDTO{
				Code: s.Code,
				Name: s.Name,
				Icon: s.Icon,
			})
		}

		var tags []string
		json.Unmarshal([]byte(doc.Tags), &tags)

		agencyName := ""
		if doc.LeadAgency != nil {
			agencyName = doc.LeadAgency.Name
		}

		createdAtVal := time.Time{}
		if doc.CreatedAt != nil {
			createdAtVal = *doc.CreatedAt
		}

		items = append(items, model.DocumentListItem{
			ID:          doc.ID,
			Title:       doc.Title,
			Slug:        doc.Slug,
			Description: doc.Description,
			Agency:      agencyName,
			Year:        doc.Year,
			Language:    doc.Language,
			FileSize:    doc.FileSize,
			TotalPages:  doc.TotalPages,
			Type:        "Report",
			PubStatus:   doc.Status,
			CoverImage:  doc.CoverImage,
			Sdgs:        sdgs,
			Tags:        tags,
			Views:       doc.Views,
			Downloads:   doc.Downloads,
			IsActive:    doc.IsActive,
			CreatedAt:   createdAtVal,
		})
	}

	page := 1
	limit := 12
	if p, ok := filters["page"].(int); ok && p > 0 {
		page = p
	}
	if l, ok := filters["limit"].(int); ok && l > 0 {
		limit = l
	}

	totalPages := 0
	if totalItems > 0 {
		totalPages = (totalItems + limit - 1) / limit
	}

	return &model.DocumentListResponse{
		Items: items,
		Pagination: model.Pagination{
			Page:       page,
			Limit:      limit,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *documentService) SearchPublicDocuments(q string, page int, limit int, sort string, filters map[string]interface{}) (map[string]interface{}, error) {
	filters["q"] = q
	filters["page"] = page
	filters["limit"] = limit
	filters["sort"] = sort

	listResp, err := s.ListPublicDocuments(filters)
	if err != nil {
		return nil, err
	}

	var highlightedItems []map[string]interface{}
	for _, item := range listResp.Items {
		hlTitle := highlightWord(item.Title, q)
		hlDesc := highlightWord(item.Description, q)

		itemMap := map[string]interface{}{
			"id":          item.ID,
			"title":       item.Title,
			"slug":        item.Slug,
			"description": item.Description,
			"agency":      item.Agency,
			"year":        item.Year,
			"language":    item.Language,
			"file_size":   item.FileSize,
			"total_pages": item.TotalPages,
			"type":        item.Type,
			"pub_status":  item.PubStatus,
			"cover_image": item.CoverImage,
			"sdgs":        item.Sdgs,
			"tags":        item.Tags,
			"views":       item.Views,
			"downloads":   item.Downloads,
			"highlight": map[string]string{
				"title":       hlTitle,
				"description": hlDesc,
			},
		}
		highlightedItems = append(highlightedItems, itemMap)
	}

	suggestions := []string{"Green Economy", "Carbon Emission", "SDGs", "Paris Agreement"}

	return map[string]interface{}{
		"success": true,
		"message": "Search results retrieved successfully",
		"data": map[string]interface{}{
			"items":       highlightedItems,
			"pagination":  listResp.Pagination,
			"suggestions": suggestions,
		},
	}, nil
}

func (s *documentService) DeleteSubmission(id string) error {
	return s.docRepo.Delete(id)
}

func (s *documentService) PublishDocument(id string) (*model.Document, error) {
	doc, err := s.docRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	doc.Status = "published"
	if err := s.docRepo.Update(doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *documentService) UnpublishDocument(id string) (*model.Document, error) {
	doc, err := s.docRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	doc.Status = "approved_unpublished"
	if err := s.docRepo.Update(doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *documentService) GetRelatedDocuments(id string) ([]model.DocumentResponse, error) {
	doc, err := s.docRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	docs, err := s.docRepo.GetRelated(doc)
	if err != nil {
		return nil, err
	}

	var resp []model.DocumentResponse
	for _, d := range docs {
		resp = append(resp, mapToDocumentResponse(&d))
	}
	return resp, nil
}

func (s *documentService) GenerateDownloadLink(id string) (map[string]interface{}, error) {
	doc, err := s.docRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Increment downloads count
	doc.Downloads += 1
	s.docRepo.Update(doc)

	expiresAt := time.Now().Add(1 * time.Hour)

	return map[string]interface{}{
		"download_url": doc.FileURL,
		"filename":     fmt.Sprintf("%s_%s.pdf", doc.Code, strings.ReplaceAll(doc.Title, " ", "_")),
		"file_size":    doc.FileSize,
		"expires_at":   expiresAt.Format(time.RFC3339),
	}, nil
}

func (s *documentService) GetPlatformStats() (map[string]interface{}, error) {
	return s.docRepo.GetGlobalStats()
}

func (s *documentService) GetAnalyticsOverview() (map[string]interface{}, error) {
	return s.docRepo.GetOverviewAnalytics()
}

func (s *documentService) GetUploadsOverTime(fromYear, toYear int) ([]map[string]interface{}, error) {
	return s.docRepo.GetUploadsOverTime(fromYear, toYear)
}

func (s *documentService) GetBySdgAnalytics() ([]map[string]interface{}, error) {
	return s.docRepo.GetBySdgAnalytics()
}

func (s *documentService) GetByAgencyAnalytics() ([]map[string]interface{}, error) {
	return s.docRepo.GetByAgencyAnalytics()
}

func (s *documentService) GetBySectorAnalytics() ([]map[string]interface{}, error) {
	return s.docRepo.GetBySectorAnalytics()
}

func (s *documentService) GetByLanguageAnalytics() ([]map[string]interface{}, error) {
	return s.docRepo.GetByLanguageAnalytics()
}

func (s *documentService) ListSubmissions(status string, search string, page int, limit int) (map[string]interface{}, error) {
	docs, totalItems, err := s.docRepo.ListSubmissions(status, search, page, limit)
	if err != nil {
		return nil, err
	}

	var items []map[string]interface{}
	for _, doc := range docs {
		agencyName := ""
		if doc.LeadAgency != nil {
			agencyName = doc.LeadAgency.Name
		}
		authorName := ""
		if doc.Author != nil && doc.Author.Name != nil {
			authorName = *doc.Author.Name
		}

		items = append(items, map[string]interface{}{
			"id":                  doc.ID,
			"title":               doc.Title,
			"short_description":   doc.Description,
			"date_of_publication": doc.DateOfPublication,
			"status":              doc.Status,
			"agency":              agencyName,
			"author":              authorName,
			"created_at":          doc.CreatedAt,
		})
	}

	totalPages := 0
	if totalItems > 0 {
		totalPages = (totalItems + limit - 1) / limit
	}

	return map[string]interface{}{
		"items": items,
		"pagination": map[string]interface{}{
			"page":       page,
			"limit":      limit,
			"totalItems": totalItems,
			"totalPages": totalPages,
		},
	}, nil
}

// Private utilities
func slugify(title string) string {
	slug := strings.ToLower(title)
	reHTML := regexp.MustCompile("<[^>]*>")
	slug = reHTML.ReplaceAllString(slug, "")
	re := regexp.MustCompile("[^a-z0-9]+")
	slug = re.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}

func generateRandomHex(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "abc"
	}
	return hex.EncodeToString(bytes)
}

func highlightWord(text, query string) string {
	if query == "" {
		return text
	}
	escapedQuery := regexp.QuoteMeta(query)
	re := regexp.MustCompile("(?i)" + escapedQuery)
	return re.ReplaceAllStringFunc(text, func(m string) string {
		return "<mark>" + m + "</mark>"
	})
}

func mapToDocumentResponse(doc *model.Document) model.DocumentResponse {
	createdAtVal := time.Time{}
	if doc.CreatedAt != nil {
		createdAtVal = *doc.CreatedAt
	}
	dateAdded := createdAtVal.Format("2006-01-02")

	var sdgs []model.SdgDTO
	for _, s := range doc.Sdgs {
		sdgs = append(sdgs, model.SdgDTO{
			Code: s.Code,
			Name: s.Name,
			Icon: s.Icon,
		})
	}

	var sectors []string
	for _, s := range doc.Sectors {
		sectors = append(sectors, s.Name)
	}

	var lnobs []string
	for _, l := range doc.Lnobs {
		lnobs = append(lnobs, l.Name)
	}

	var tags []string
	var thematicAreas []string
	var otherAgencies []string
	var partners []model.PartnerDTO
	var supportingFiles []model.SupportingFile

	json.Unmarshal([]byte(doc.Tags), &tags)
	json.Unmarshal([]byte(doc.ThematicAreas), &thematicAreas)
	json.Unmarshal([]byte(doc.OtherAgencies), &otherAgencies)
	json.Unmarshal([]byte(doc.NonUnPartners), &partners)
	json.Unmarshal([]byte(doc.SupportingFiles), &supportingFiles)

	leadAgencyName := ""
	if doc.LeadAgency != nil {
		leadAgencyName = doc.LeadAgency.Name
	} else {
		leadAgencyName = doc.LeadAgencyCode
	}

	jpName := ""
	if doc.JointProgramme != nil {
		jpName = doc.JointProgramme.Name
	} else {
		jpName = doc.JointProgrammeCode
	}

	updatedAtVal := time.Time{}
	if doc.UpdatedAt != nil {
		updatedAtVal = *doc.UpdatedAt
	}

	return model.DocumentResponse{
		ID:         doc.ID,
		Code:       doc.Code,
		Slug:       doc.Slug,
		Title:      doc.Title,
		Agency:     leadAgencyName,
		Year:       doc.Year,
		Language:   doc.Language,
		FileURL:    doc.FileURL,
		FileSize:   doc.FileSize,
		DateAdded:  dateAdded,
		Type:       "Report",
		TotalPages: doc.TotalPages,
		PubStatus:  doc.Status,
		CoverImage: doc.CoverImage,
		Abstract:   doc.Abstract,
		Summary:    doc.Summary,
		Sdgs:       sdgs,
		Tags:       tags,
		ThematicAreas: thematicAreas,
		Sectors:    sectors,
		LnobGroups: lnobs,
		Classification: model.Classification{
			LeadAgency:      leadAgencyName,
			OtherAgencies:   otherAgencies,
			JointProgramme:  jpName,
			GeographicScope: doc.GeographicScope,
			NonUnPartners:   partners,
		},
		FocalPoint: model.FocalPointDTO{
			Name:       doc.FocalPointName,
			Email:      doc.FocalPointEmail,
			Phone:      doc.FocalPointPhone,
			Department: doc.FocalPointDepartment,
		},
		Views:           doc.Views,
		Downloads:       doc.Downloads,
		CreatedAt:       createdAtVal,
		UpdatedAt:       updatedAtVal,
		IsActive:        doc.IsActive,
		SupportingFiles: supportingFiles,
	}
}

func (s *documentService) RecordDocumentActivity(docID string, action string, ip string, userAgent string) error {
	logEntry := &model.DocumentActivityLog{
		DocumentID: docID,
		Action:     action,
		IPAddress:  ip,
		UserAgent:  userAgent,
	}
	return s.docRepo.RecordActivity(logEntry)
}
