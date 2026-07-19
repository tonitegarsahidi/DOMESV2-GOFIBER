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
	app.Static("/public/upload", "./public/upload")

	// Repositories
	userRepo := repository.NewUserRepository()
	healthRepo := repository.NewHealthRepository()
	masterRepo := repository.NewMasterRepository()
	docRepo := repository.NewDocumentRepository()
	reportRepo := repository.NewReportRepository()
	cmsRepo := repository.NewCmsRepository()

	// Services
	mailService := service.NewMailService()
	authService := service.NewAuthService(userRepo, mailService)
	healthService := service.NewHealthService(healthRepo)
	masterService := service.NewMasterService(masterRepo)
	docService := service.NewDocumentService(docRepo, userRepo)
	reportService := service.NewReportService(reportRepo)
	cmsService := service.NewCmsService(cmsRepo, userRepo)
	uploadService := service.NewFileUploadService()

	// Controllers
	authController := controller.NewAuthController(authService)
	healthController := controller.NewHealthController(healthService)
	masterController := controller.NewMasterController(masterService)
	uploadController := controller.NewUploadController(uploadService)
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

		// Public Master Endpoints
		master := api.Group("/master")
		{
			master.Get("/agencies", masterController.GetAgencies)
			master.Get("/sdgs", masterController.GetSdgs)
			master.Get("/sectors", masterController.GetSectors)
			master.Get("/languages", masterController.GetLanguages)
			master.Get("/joint-programmes", masterController.GetJointProgrammes)
			master.Get("/lnobs", masterController.GetLnobs)
			master.Get("/non-un-partners", masterController.GetNonUnPartners)
			master.Get("/organizations", masterController.GetOrganizations)
			master.Get("/thematic-areas", masterController.GetThematicAreas)
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
			cms := protected.Group("/cms")
			{
				// User Profiles & Settings
				user := cms.Group("/user")
				{
					user.Get("/me", authController.Me)
					user.Put("/profile", authController.UpdateProfile)
					user.Put("/password", authController.ChangePassword)
					user.Get("/notifications", authController.GetNotificationPreferences)
					user.Put("/notifications", authController.UpdateNotificationPreferences)
				}

				// Admin whitelist settings
				admin := cms.Group("/admin")
				{
					admin.Get("/emails", authController.GetAdminEmails)
					admin.Post("/emails", authController.AddAdminEmail)
					admin.Delete("/emails/:email", authController.DeleteAdminEmail)
				}

				// CMS Dashboard
				cms.Get("/dashboard", cmsController.GetDashboardStats)
				cms.Get("/activity", cmsController.GetRecentActivity)

				// CMS Master Management
				cms.Get("/master/:type", cmsController.ListMasters)
				cms.Post("/master/:type", cmsController.CreateMaster)
				cms.Put("/master/:type/:code", cmsController.UpdateMaster)
				cms.Delete("/master/:type/:code", cmsController.DeleteMaster)

				// CMS Submissions Wizard & Mgmt
				cms.Get("/submissions", docController.ListSubmissions)
				cms.Post("/submissions", docController.CreateSubmission)
				cms.Put("/submissions/:id", docController.UpdateSubmission)
				cms.Post("/submissions/:id/draft", docController.SaveDraft)
				cms.Delete("/submissions/:id", docController.DeleteSubmission)
				cms.Put("/submissions/:id/publish", docController.PublishDocument)
				cms.Put("/submissions/:id/unpublish", docController.UnpublishDocument)

				// CMS Users Management (Admin only verified in Controller)
				cms.Get("/users", cmsController.ListUsers)
				cms.Post("/users", cmsController.CreateUser)
				cms.Put("/users/:id", cmsController.UpdateUser)
				cms.Delete("/users/:id", cmsController.DeleteUser)

				// CMS Reports Management
				cms.Get("/reports", reportController.ListReports)
				cms.Put("/reports/:id/status", reportController.UpdateStatus)

				// CMS Analytics
				cms.Get("/analytics/summary", cmsController.GetAnalyticsSummary)
				cms.Get("/analytics/top-downloads", cmsController.GetTopDownloads)
				cms.Get("/analytics/top-views", cmsController.GetTopViews)

				// File Upload Protected Actions
				cms.Post("/upload", uploadController.UploadFile)
				cms.Post("/upload/url-validate", uploadController.ValidateURL)
				cms.Post("/upload/avatar", uploadController.UploadAvatar)
			}
		}
	}
}
