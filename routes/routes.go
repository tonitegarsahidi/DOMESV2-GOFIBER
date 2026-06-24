package routes

import (
	"domesv2/internal/controller"
	"domesv2/internal/middleware"
	"domesv2/internal/repository"
	"domesv2/internal/service"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Static files serving for uploaded assets
	app.Static("/uploads", "./uploads")

	// Repositories
	userRepo := repository.NewUserRepository()
	healthRepo := repository.NewHealthRepository()
	refRepo := repository.NewReferenceRepository()
	docRepo := repository.NewDocumentRepository()
	reportRepo := repository.NewReportRepository()
	cmsRepo := repository.NewCmsRepository()

	// Services
	mailService := service.NewMailService()
	authService := service.NewAuthService(userRepo, mailService)
	healthService := service.NewHealthService(healthRepo)
	refService := service.NewReferenceService(refRepo)
	docService := service.NewDocumentService(docRepo, userRepo)
	reportService := service.NewReportService(reportRepo)
	cmsService := service.NewCmsService(cmsRepo, userRepo)

	// Controllers
	authController := controller.NewAuthController(authService)
	healthController := controller.NewHealthController(healthService)
	refController := controller.NewReferenceController(refService)
	uploadController := controller.NewUploadController()
	docController := controller.NewDocumentController(docService)
	reportController := controller.NewReportController(reportService)
	cmsController := controller.NewCmsController(cmsService, authService)

	// Base Route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Hello Domes v2",
			"version": "1.0.0",
			"docs":    "/api/v2/health-check",
		})
	})

	api := app.Group("/api/v2")
	{
		// Health Check
		api.Get("/health-check", healthController.Check)

		// Authentication Public Endpoints
		auth := api.Group("/auth")
		{
			auth.Post("/register", authController.Register)
			auth.Post("/login", authController.Login)
			auth.Post("/forgot-password", authController.ForgotPassword)
			auth.Post("/reset-password", authController.ResetPassword)
		}

		// Public Reference Endpoints
		ref := api.Group("/reference")
		{
			ref.Get("/agencies", refController.GetAgencies)
			ref.Get("/sdgs", refController.GetSdgs)
			ref.Get("/sectors", refController.GetSectors)
			ref.Get("/languages", refController.GetLanguages)
			ref.Get("/joint-programmes", refController.GetJointProgrammes)
			ref.Get("/lnobs", refController.GetLnobs)
			ref.Get("/non-un-partners", refController.GetNonUnPartners)
			ref.Get("/organizations", refController.GetOrganizations)
		}

		// Public Documents Endpoints
		api.Get("/documents", docController.ListPublic)
		api.Get("/documents/search", docController.SearchPublic)
		api.Get("/documents/:id", docController.GetByIDOrSlug)
		api.Get("/documents/:id/related", docController.GetRelated)
		api.Get("/documents/:id/download", docController.Download)

		// Public Stats & Public Analytics
		api.Get("/stats", docController.GetPlatformStats)
		api.Get("/analytics/overview", docController.GetAnalyticsOverview)
		api.Get("/analytics/uploads-over-time", docController.GetUploadsOverTime)
		api.Get("/analytics/by-sdg", docController.GetBySdgAnalytics)
		api.Get("/analytics/by-agency", docController.GetByAgencyAnalytics)
		api.Get("/analytics/by-sector", docController.GetBySectorAnalytics)
		api.Get("/analytics/by-language", docController.GetByLanguageAnalytics)

		// Public Broken Link Report
		api.Post("/reports", reportController.SubmitReport)

		// Protected Routes Group
		protected := api.Group("/")
		protected.Use(middleware.JWTMiddleware())
		{
			// User Profiles & Settings
			user := protected.Group("/user")
			{
				user.Get("/me", authController.Me)
				user.Put("/profile", authController.UpdateProfile)
				user.Put("/password", authController.ChangePassword)
				user.Get("/notifications", authController.GetNotificationPreferences)
				user.Put("/notifications", authController.UpdateNotificationPreferences)
			}

			// Admin whitelist settings
			admin := protected.Group("/admin")
			{
				admin.Get("/emails", authController.GetAdminEmails)
				admin.Post("/emails", authController.AddAdminEmail)
				admin.Delete("/emails/:email", authController.DeleteAdminEmail)
			}

			// CMS Dashboard
			protected.Get("/cms/dashboard", cmsController.GetDashboardStats)
			protected.Get("/cms/activity", cmsController.GetRecentActivity)

			// CMS Reference Management
			protected.Get("/cms/reference/:type", cmsController.ListReferences)
			protected.Post("/cms/reference/:type", cmsController.CreateReference)
			protected.Put("/cms/reference/:type/:code", cmsController.UpdateReference)
			protected.Delete("/cms/reference/:type/:code", cmsController.DeleteReference)

			// CMS Submissions Wizard & Mgmt
			protected.Get("/submissions", docController.ListSubmissions)
			protected.Post("/submissions", docController.CreateSubmission)
			protected.Post("/submissions/:id/draft", docController.SaveDraft)
			protected.Delete("/submissions/:id", docController.DeleteSubmission)
			protected.Put("/submissions/:id/publish", docController.PublishDocument)
			protected.Put("/submissions/:id/unpublish", docController.UnpublishDocument)

			// CMS Users Management (Admin only verified in Controller)
			protected.Get("/users", cmsController.ListUsers)
			protected.Post("/users", cmsController.CreateUser)
			protected.Put("/users/:id", cmsController.UpdateUser)
			protected.Delete("/users/:id", cmsController.DeleteUser)

			// CMS Reports Management
			protected.Get("/reports", reportController.ListReports)
			protected.Put("/reports/:id/status", reportController.UpdateStatus)

			// CMS Analytics
			protected.Get("/analytics/summary", cmsController.GetAnalyticsSummary)
			protected.Get("/analytics/top-downloads", cmsController.GetTopDownloads)
			protected.Get("/analytics/top-views", cmsController.GetTopViews)

			// File Upload Protected Actions
			protected.Post("/upload", uploadController.UploadFile)
			protected.Post("/upload/url-validate", uploadController.ValidateURL)
			protected.Post("/upload/avatar", uploadController.UploadAvatar)
		}
	}
}
