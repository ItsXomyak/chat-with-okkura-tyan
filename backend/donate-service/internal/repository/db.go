package repository

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"donate-service/internal/model"
)

func NewPostgresDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&model.Donation{})
	db.AutoMigrate(&model.UserProfile{})

	return db, nil
}
