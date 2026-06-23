package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"domesv2/config"
	"domesv2/config/database"
	"domesv2/config/logger"
	"domesv2/config/redis"
	"domesv2/internal/middleware"
	"domesv2/routes"
	"go.uber.org/zap"
)

func main() {
	// Initialize configuration
	config.InitConfig()
	cfg := config.AppConfig

	// Initialize logger
	logger.InitLogger(cfg.Server.Env)
	defer logger.Sync()

	// Initialize database
	database.InitMySQL(cfg)

	// Initialize Redis (optional)
	redis.InitRedis(cfg)

	// Auto migrate models (uncomment when you have models to migrate)
	// db := database.GetDB()
	// db.AutoMigrate(&model.User{})

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "DOMESv2 API",
		ErrorHandler: middleware.GlobalErrorHandler,
	})

	// Global middlewares
	app.Use(recover.New())
	app.Use(middleware.LoggingMiddleware())

	// Setup routes
	routes.SetupRoutes(app)

	// Start server
	logger.Info("Server starting", zap.String("port", cfg.Server.Port))
	if err := app.Listen(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
