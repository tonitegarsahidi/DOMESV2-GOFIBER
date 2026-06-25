package main

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	database.MigrateAndSeed(database.GetDB())

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
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			// Allow all origins in local or development environments
			env := strings.ToLower(cfg.Server.Env)
			if env == "local" || env == "development" || env == "dev" || env == "" {
				return true
			}

			// If allowed origins are explicitly specified in .env, check them
			if cfg.Server.AllowedOrigins != "" {
				origins := strings.Split(cfg.Server.AllowedOrigins, ",")
				for _, o := range origins {
					o = strings.TrimSpace(o)
					if o == "*" || o == origin {
						return true
					}
				}
			}

			// Allow localhost and 127.0.0.1 origins for local integration/testing
			if strings.Contains(origin, "://localhost") || strings.Contains(origin, "://127.0.0.1") {
				return true
			}

			return false
		},
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		AllowMethods:     "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS",
		AllowCredentials: true,
	}))
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
