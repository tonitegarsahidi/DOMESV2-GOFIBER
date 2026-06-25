package main

import (
	"domesv2/config"
	"domesv2/config/database"
	"domesv2/config/logger"
	"log"
)

func main() {
	// Initialize configuration
	config.InitConfig()
	cfg := config.AppConfig

	// Initialize logger
	logger.InitLogger(cfg.Server.Env)

	// Initialize database connection
	database.InitMySQL(cfg)
	db := database.GetDB()
	if db == nil {
		log.Fatal("Could not establish database connection")
	}

	// Run migration and seeders
	log.Println("Starting database migrations and seeders...")
	database.MigrateAndSeed(db)
	log.Println("Database migrations and seeders completed successfully.")
}
