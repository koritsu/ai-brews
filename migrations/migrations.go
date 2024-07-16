package migrations

import (
	"ai-brews/config"
	"ai-brews/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
)

func Migrate() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config file")
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = cfg.DBPath
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database")
	}

	db.AutoMigrate(&models.Recipe{})
	db.AutoMigrate(&models.Admin{})
}
