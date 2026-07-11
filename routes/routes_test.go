package routes_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"

	"domesv2/config"
	"domesv2/config/database"
	"domesv2/config/logger"
	"domesv2/internal/middleware"
	"domesv2/internal/model"
	"domesv2/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

var app *fiber.App
var token string
var testDocID string
var testDocSlug string
var testReportID string
var testUserID uint

func TestMain(m *testing.M) {
	// Load env from parent directory (since tests run in the package dir)
	_ = godotenv.Load("../.env")

	config.InitConfig()
	logger.InitLogger("development")
	database.InitMySQL(config.AppConfig)
	db := database.GetDB()
	if db != nil {
		db.Exec("DROP TABLE IF EXISTS V2Reports")
		db.Exec("DROP TABLE IF EXISTS v2_document_sdgs")
		db.Exec("DROP TABLE IF EXISTS v2_document_sectors")
		db.Exec("DROP TABLE IF EXISTS v2_document_lnobs")
		db.Exec("DROP TABLE IF EXISTS V2DocumentStats")
		db.Exec("DROP TABLE IF EXISTS V2DocumentActivityLogs")
		db.Exec("DROP TABLE IF EXISTS V2Documents")
	}
	database.MigrateAndSeed(db)

	// Set user tonitegarsahidi@gmail.com as administrator for testing protected admin endpoints
	if db != nil {
		db.Model(&model.User{}).Where("email = ?", "tonitegarsahidi@gmail.com").Updates(map[string]interface{}{
			"role":   "administrator",
			"type":   "admin",
			"status": "active",
		})
		// Clean up previous test submissions/reports to prevent duplicate entry/validation failures
		db.Exec("DELETE FROM V2Reports")
		db.Exec("DELETE FROM v2_document_sdgs")
		db.Exec("DELETE FROM v2_document_sectors")
		db.Exec("DELETE FROM v2_document_lnobs")
		db.Exec("DELETE FROM V2Documents")
		db.Exec("DELETE FROM V2MasterSectors WHERE code = 'test-sector-xyz'")
		db.Exec("DELETE FROM V2MasterThematicAreas WHERE code = 'test-sector-xyz'")
		db.Exec("DELETE FROM V2AdminEmails WHERE email = 'test-admin@un.org'")
	}

	app = fiber.New(fiber.Config{
		ErrorHandler: middleware.GlobalErrorHandler,
	})
	app.Use(recover.New())
	app.Use(middleware.LoggingMiddleware())
	routes.SetupRoutes(app)

	m.Run()
}

func TestA_PublicBaseAndHealth(t *testing.T) {
	// GET /
	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test GET /: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	// GET /api/v2/health-check
	req = httptest.NewRequest("GET", "/api/v2/health-check", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test GET /api/v2/health-check: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestB_AuthLogin(t *testing.T) {
	loginPayload := model.LoginRequest{
		Email:    "tonitegarsahidi@gmail.com",
		Password: "rahasiaku123",
		Captcha:  "RANDOM",
	}
	body, _ := json.Marshal(loginPayload)
	req := httptest.NewRequest("POST", "/api/v2/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test POST /api/v2/auth/login: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected 200, got %d. Body: %s", resp.StatusCode, string(bodyBytes))
	}

	var res map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&res)

	success := res["success"].(bool)
	if !success {
		t.Errorf("Expected success true, got false")
	}

	data := res["data"].(map[string]interface{})
	token = data["token"].(string)
	if token == "" {
		t.Errorf("Expected non-empty token")
	}
}

func TestC_UserProfileAndSettings(t *testing.T) {
	if token == "" {
		t.Skip("Token is empty, skipping protected profile tests")
	}

	// GET /api/v2/user/me
	req := httptest.NewRequest("GET", "/api/v2/user/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test GET /api/v2/user/me: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	// PUT /api/v2/user/profile
	profilePayload := model.UpdateProfileRequest{
		FirstName:    "Toni",
		LastName:     "Tegar",
		Position:     "Senior Administrator",
		Organization: "UNDP",
		PhoneNumber:  "+628123456789",
	}
	body, _ := json.Marshal(profilePayload)
	req = httptest.NewRequest("PUT", "/api/v2/user/profile", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test PUT /api/v2/user/profile: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	// GET /api/v2/user/notifications
	req = httptest.NewRequest("GET", "/api/v2/user/notifications", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test GET /api/v2/user/notifications: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	// PUT /api/v2/user/notifications
	truVal := true
	falVal := false
	notifPayload := model.UpdateNotificationRequest{
		DocumentApprovals:  &truVal,
		BrokenLinkReports:  &falVal,
		SystemUpdates:      &truVal,
		EmailNotifications: &truVal,
	}
	body, _ = json.Marshal(notifPayload)
	req = httptest.NewRequest("PUT", "/api/v2/user/notifications", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test PUT /api/v2/user/notifications: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestD_AdminEmailsWhitelist(t *testing.T) {
	if token == "" {
		t.Skip("Token is empty, skipping whitelist tests")
	}

	// POST /api/v2/admin/emails
	emailPayload := model.AddAdminEmailRequest{
		Email: "test-admin@un.org",
	}
	body, _ := json.Marshal(emailPayload)
	req := httptest.NewRequest("POST", "/api/v2/admin/emails", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test POST /api/v2/admin/emails: %v", err)
	}
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusConflict {
		t.Errorf("Expected 201 or 409, got %d", resp.StatusCode)
	}

	// GET /api/v2/admin/emails
	req = httptest.NewRequest("GET", "/api/v2/admin/emails", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test GET /api/v2/admin/emails: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	// DELETE /api/v2/admin/emails/test-admin@un.org
	req = httptest.NewRequest("DELETE", "/api/v2/admin/emails/test-admin@un.org", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test DELETE /api/v2/admin/emails/:email: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestE_PublicMasterData(t *testing.T) {
	masters := []string{
		"agencies", "sdgs", "sectors", "languages", "joint-programmes", "lnobs", "non-un-partners", "organizations", "thematic-areas",
	}

	for _, endpoint := range masters {
		t.Run("GET /api/v2/master/"+endpoint, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v2/master/"+endpoint, nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to query /api/v2/master/%s: %v", endpoint, err)
			}
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected 200, got %d", resp.StatusCode)
			}
		})
	}
}

func TestF_PublicStatsAndAnalytics(t *testing.T) {
	analytics := []string{
		"overview", "uploads-over-time", "by-sdg", "by-agency", "by-sector", "by-language",
	}

	// GET /api/v2/stats
	req := httptest.NewRequest("GET", "/api/v2/stats", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed GET /api/v2/stats: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	for _, endpoint := range analytics {
		t.Run("GET /api/v2/analytics/"+endpoint, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v2/analytics/"+endpoint, nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed GET /api/v2/analytics/%s: %v", endpoint, err)
			}
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected 200, got %d", resp.StatusCode)
			}
		})
	}
}

func TestG_CmsDashboardAndActivity(t *testing.T) {
	if token == "" {
		t.Skip("Token is empty, skipping CMS Dashboard tests")
	}

	// GET /api/v2/cms/dashboard
	req := httptest.NewRequest("GET", "/api/v2/cms/dashboard", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed GET /api/v2/cms/dashboard: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	// GET /api/v2/cms/activity
	req = httptest.NewRequest("GET", "/api/v2/cms/activity", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed GET /api/v2/cms/activity: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestH_SubmissionsWizardAndPublishing(t *testing.T) {
	if token == "" {
		t.Skip("Token is empty, skipping submissions tests")
	}

	// 1. POST /api/v2/submissions/:id/draft (Create first draft step 2)
	draftPayload := model.DraftRequest{
		Step: 2,
		Data: map[string]interface{}{
			"title":                  "Digital Economy and Financial Inclusion in Rural Indonesia",
			"date_of_publication":    "2024-06-15",
			"total_pages":            120,
			"language":               "English, Bahasa Indonesia",
			"publication_status":     "Published",
			"short_summary":          "Analysis of digital financial services...",
			"tags":                   []string{"digital economy", "financial inclusion", "fintech"},
			"focal_point_name":       "Budi Santoso",
			"focal_point_email":      "b.santoso@undp.org",
			"focal_point_phone":      "+62 812 3456 7890",
			"focal_point_department": "Inclusive Growth Unit",
		},
	}
	body, _ := json.Marshal(draftPayload)
	req := httptest.NewRequest("POST", "/api/v2/submissions/0/draft", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed draft submission: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	var draftRes map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&draftRes)
	draftData := draftRes["data"].(map[string]interface{})
	testDocID = draftData["id"].(string)

	// 2. POST /api/v2/submissions (Final Submit step 4)
	submitPayload := model.SubmissionRequest{
		Title:             "Digital Economy and Financial Inclusion in Rural Indonesia",
		ShortDescription:  "Analysis of digital financial services...",
		Abstract:          "This comprehensive report analyzes...",
		DetailedSummary:   "<b>Executive Overview</b><br><br>This extensive report...",
		DateOfPublication: "2024-06-15",
		TotalPages:        120,
		Language:          "English, Bahasa Indonesia",
		PublicationStatus: "Published",
		Tags:              []string{"digital economy", "financial inclusion", "fintech"},
		FileURL:           "/uploads/doc_test.pdf",
		FileSize:          "4.2 MB",
		CoverImageURL:     "/uploads/cover_test.jpg",
		ExternalURL:       "https://un.or.id/reports/digital-economy",
		SupportingFiles: []model.SupportingFile{
			{URL: "/uploads/infographic.pdf", Type: "infographic", Description: "Infographic details"},
		},
		Agency: "UNDP",
		FocalPoint: model.FocalPointDTO{
			Name:       "Budi Santoso",
			Email:      "b.santoso@undp.org",
			Phone:      "+62 812 3456 7890",
			Department: "Inclusive Growth Unit",
		},
		Sdgs:            []string{"GOAL 1", "GOAL 8", "GOAL 9"},
		Sectors:         []string{"economic-development", "innovation-technology"},
		LnobGroups:      []string{"women-girls", "others"},
		JointProgramme:  "assisst",
		OtherAgencies:   []string{"UN Women", "UNICEF"},
		NonUnPartners:   []model.PartnerDTO{{Type: "government", Name: "Ministry of Finance"}},
		ThematicAreas:   []string{"Digital Inclusion", "Economic Growth"},
		GeographicScope: "National",
	}
	body, _ = json.Marshal(submitPayload)
	req = httptest.NewRequest("POST", "/api/v2/submissions", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed final submit: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected 201, got %d. Body: %s", resp.StatusCode, string(bodyBytes))
	}

	var submitRes map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&submitRes)
	submitData := submitRes["data"].(map[string]interface{})
	testDocID = submitData["id"].(string)

	// 3. PUT /api/v2/submissions/:id/publish
	req = httptest.NewRequest("PUT", fmt.Sprintf("/api/v2/submissions/%s/publish", testDocID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed publish document: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	// Fetch doc slug for subsequent testing
	db := database.GetDB()
	var doc model.Document
	db.First(&doc, "id = ?", testDocID)
	testDocSlug = doc.Slug

	// 4. GET /api/v2/submissions (List submissions CMS)
	req = httptest.NewRequest("GET", "/api/v2/submissions?status=all", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed list submissions: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	// 5. PUT /api/v2/submissions/:id/unpublish
	req = httptest.NewRequest("PUT", fmt.Sprintf("/api/v2/submissions/%s/unpublish", testDocID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed unpublish document: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	// Publish again so public documents tests can query it
	req = httptest.NewRequest("PUT", fmt.Sprintf("/api/v2/submissions/%s/publish", testDocID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	_, _ = app.Test(req)
}

func TestI_PublicDocumentsQueries(t *testing.T) {
	// GET /api/v2/documents
	req := httptest.NewRequest("GET", "/api/v2/documents", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed GET /api/v2/documents: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	// GET /api/v2/documents/search
	req = httptest.NewRequest("GET", "/api/v2/documents/search?q=Digital", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed GET /api/v2/documents/search: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	if testDocID != "" {
		// GET /api/v2/documents/{id}
		req = httptest.NewRequest("GET", fmt.Sprintf("/api/v2/documents/%s", testDocID), nil)
		resp, err = app.Test(req)
		if err != nil {
			t.Fatalf("Failed GET /api/v2/documents/{id}: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}

		// GET /api/v2/documents/{slug}
		req = httptest.NewRequest("GET", "/api/v2/documents/"+testDocSlug, nil)
		resp, err = app.Test(req)
		if err != nil {
			t.Fatalf("Failed GET /api/v2/documents/{slug}: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}

		// GET /api/v2/documents/{id}/related
		req = httptest.NewRequest("GET", fmt.Sprintf("/api/v2/documents/%s/related", testDocID), nil)
		resp, err = app.Test(req)
		if err != nil {
			t.Fatalf("Failed GET /api/v2/documents/{id}/related: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}

		// GET /api/v2/documents/{id}/download
		req = httptest.NewRequest("GET", fmt.Sprintf("/api/v2/documents/%s/download", testDocID), nil)
		resp, err = app.Test(req)
		if err != nil {
			t.Fatalf("Failed GET /api/v2/documents/{id}/download: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	}
}

func TestJ_ReportsHandling(t *testing.T) {
	// 1. POST /api/v2/reports (Public submit broken link)
	reportPayload := model.CreateReportRequest{
		DocumentID:    testDocID,
		ReporterName:  "Budi Santoso",
		ReporterEmail: "budi@example.com",
		Details:       "The PDF link leads to a 404 error page.",
	}
	body, _ := json.Marshal(reportPayload)
	req := httptest.NewRequest("POST", "/api/v2/reports", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed public submit report: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected 201, got %d", resp.StatusCode)
	}

	var res map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&res)
	reportData := res["data"].(map[string]interface{})
	testReportID = reportData["id"].(string)

	if token == "" {
		t.Skip("Token is empty, skipping CMS reports check")
	}

	// 2. GET /api/v2/reports (CMS list reports)
	req = httptest.NewRequest("GET", "/api/v2/reports?status=all", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed CMS list reports: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	// 3. PUT /api/v2/reports/:id/status (CMS update status)
	statusPayload := model.UpdateReportStatusRequest{
		Status: "in_progress",
	}
	body, _ = json.Marshal(statusPayload)
	req = httptest.NewRequest("PUT", fmt.Sprintf("/api/v2/reports/%s/status", testReportID), bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed update report status: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestK_UploadsActions(t *testing.T) {
	if token == "" {
		t.Skip("Token is empty, skipping uploads actions")
	}

	// 1. POST /api/v2/upload/url-validate
	urlPayload := map[string]string{
		"url": "https://www.google.com",
	}
	body, _ := json.Marshal(urlPayload)
	req := httptest.NewRequest("POST", "/api/v2/upload/url-validate", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed URL validation: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	// 2. POST /api/v2/upload (Multipart upload file)
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="test.pdf"`)
	h.Set("Content-Type", "application/pdf")
	part, _ := w.CreatePart(h)
	_, _ = part.Write([]byte("%PDF-1.4 mock PDF content"))
	_ = w.WriteField("type", "document")
	w.Close()

	req = httptest.NewRequest("POST", "/api/v2/upload", &b)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed document upload: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected 201, got %d. Body: %s", resp.StatusCode, string(bodyBytes))
	}
}

func TestL_CmsUsersManagement(t *testing.T) {
	if token == "" {
		t.Skip("Token is empty, skipping CMS Users Management")
	}

	// 1. GET /api/v2/users
	req := httptest.NewRequest("GET", "/api/v2/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed CMS GET /api/v2/users: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	// 2. POST /api/v2/users (Create user)
	userPayload := model.CreateUserRequest{
		FirstName:       "Budi",
		LastName:        "Santoso",
		Email:           "budi.santoso@un.org",
		Password:        "password123",
		ConfirmPassword: "password123",
		Organization:    "WHO",
		Position:        "Health Officer",
		PhoneNumber:     "+6281122334455",
		Role:            "editor",
		Status:          "active",
	}
	body, _ := json.Marshal(userPayload)
	req = httptest.NewRequest("POST", "/api/v2/users", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed CMS create user: %v", err)
	}
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusConflict {
		t.Fatalf("Expected 201 or 409, got %d", resp.StatusCode)
	}

	if resp.StatusCode == http.StatusCreated {
		var userRes map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&userRes)
		userData := userRes["data"].(map[string]interface{})
		testUserID = uint(userData["id"].(float64))

		// 3. PUT /api/v2/users/{id} (Update user)
		updatePayload := model.UpdateUserRequest{
			Position: func() *string { s := "Senior Health Officer"; return &s }(),
		}
		body, _ = json.Marshal(updatePayload)
		req = httptest.NewRequest("PUT", fmt.Sprintf("/api/v2/users/%d", testUserID), bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		resp, err = app.Test(req)
		if err != nil {
			t.Fatalf("Failed CMS update user: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}

		// 4. DELETE /api/v2/users/{id} (Delete user)
		req = httptest.NewRequest("DELETE", fmt.Sprintf("/api/v2/users/%d", testUserID), nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err = app.Test(req)
		if err != nil {
			t.Fatalf("Failed CMS delete user: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	}
}

func TestM_CmsAnalyticsSummary(t *testing.T) {
	if token == "" {
		t.Skip("Token is empty, skipping CMS analytics tests")
	}

	analytics := []string{"summary", "top-downloads", "top-views"}
	for _, endpoint := range analytics {
		t.Run("GET /api/v2/analytics/"+endpoint, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v2/analytics/"+endpoint, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed CMS /api/v2/analytics/%s: %v", endpoint, err)
			}
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected 200, got %d", resp.StatusCode)
			}
		})
	}
}

func TestN_CmsMasterManagement(t *testing.T) {
	if token == "" {
		t.Skip("Token is empty, skipping CMS master management tests")
	}

	// 1. GET /api/v2/cms/master/sectors
	req := httptest.NewRequest("GET", "/api/v2/cms/master/sectors", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to GET cms master: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	// 2. POST /api/v2/cms/master/sectors (Create new sector)
	sectorPayload := model.MasterRequest{
		Code: "test-sector-xyz",
		Name: "Test Sector XYZ",
	}
	body, _ := json.Marshal(sectorPayload)
	req = httptest.NewRequest("POST", "/api/v2/cms/master/sectors", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to POST cms master: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected 201, got %d. Body: %s", resp.StatusCode, string(bodyBytes))
	}

	// 3. PUT /api/v2/cms/master/sectors/test-sector-xyz (Update sector)
	updatePayload := model.MasterRequest{
		Name: "Updated Test Sector XYZ",
	}
	body, _ = json.Marshal(updatePayload)
	req = httptest.NewRequest("PUT", "/api/v2/cms/master/sectors/test-sector-xyz", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to PUT cms master: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	// 4. DELETE /api/v2/cms/master/sectors/test-sector-xyz (Delete sector)
	req = httptest.NewRequest("DELETE", "/api/v2/cms/master/sectors/test-sector-xyz", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to DELETE cms master: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestZ_CleanupSubmissions(t *testing.T) {
	if token == "" || testDocID == "" {
		t.Skip("Skipping cleanup")
	}

	// DELETE /api/v2/submissions/{id}
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v2/submissions/%s", testDocID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed cleanup submission: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}
