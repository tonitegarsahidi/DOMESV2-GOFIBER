package routes

import (
	"domesv2/internal/controller"
	"domesv2/internal/middleware"
	"domesv2/internal/repository"
	"domesv2/internal/service"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	userRepo := repository.NewUserRepository()
	healthRepo := repository.NewHealthRepository()
	mailService := service.NewMailService()

	authService := service.NewAuthService(userRepo, mailService)
	healthService := service.NewHealthService(healthRepo)

	authController := controller.NewAuthController(authService)
	healthController := controller.NewHealthController(healthService)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Hello Domes v2",
			"version": "1.0.0",
			"docs":    "/api/health-check",
		})
	})

	api := app.Group("/api")
	{
		api.Get("/health-check", healthController.Check)

		auth := api.Group("/auth")
		{
			auth.Post("/register", authController.Register)
			auth.Post("/login", authController.Login)
			auth.Post("/forgot-password", authController.ForgotPassword)
			auth.Post("/reset-password", authController.ResetPassword)
		}

		protected := api.Group("/")
		protected.Use(middleware.JWTMiddleware())
		{
			user := protected.Group("/user")
			{
				user.Get("/me", authController.Me)
			}
		}
	}
}
