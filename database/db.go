package database

import (
	"fmt"
	"log"
	"strconv"

	"github.com/mich31/scoreplay-media-api/config"
	"github.com/mich31/scoreplay-media-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	port, err := strconv.Atoi(config.Config("DB_PORT"))
	if err != nil {
		return nil, fmt.Errorf("invalid port format: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Config("DB_HOST"), port, config.Config("DB_USER"), config.Config("DB_PASSWORD"), config.Config("DB_NAME"))
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Migrate the models
	if err := db.AutoMigrate(&models.Tag{}, &models.Media{}, &models.MediaTag{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	log.Println("Successfully connected to database")
	return db, nil
}
