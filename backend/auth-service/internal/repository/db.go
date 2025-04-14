package repository

import (
	"database/sql"
	"fmt"
	"time"

	"authforge/config"
	"authforge/internal/logger"

	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Error("Error opening DB: ", err)
		return nil, fmt.Errorf("error opening DB: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		logger.Error("Error pinging DB: ", err)
		return nil, fmt.Errorf("error pinging DB: %w", err)
	}

	return db, nil
}
