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
		log.Fatal("PORT not entered in INT format")
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Config("DB_HOST"), port, config.Config("DB_USER"), config.Config("DB_PASSWORD"), config.Config("DB_NAME"))
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the models
	db.AutoMigrate(&models.Tag{})

	fmt.Println("Connection Opened to Database")
	return db, nil
}
