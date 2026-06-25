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

	// Run legacy data migration
	log.Println("Starting legacy data migration...")
	database.MigrateLegacyData(db)
	log.Println("Legacy data migration completed successfully.")
}
