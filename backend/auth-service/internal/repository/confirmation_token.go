package repository

import (
	"database/sql"
	"time"

	"authforge/internal/logger"
	"authforge/internal/models"
)

type ConfirmationTokenRepository interface {
	CreateToken(token *models.ConfirmationToken) error
	GetTokenByString(token string) (*models.ConfirmationToken, error)
	DeleteToken(token string) error
}

type PostgresConfirmationTokenRepository struct {
	DB *sql.DB
}

func NewConfirmationTokenRepository(db *sql.DB) ConfirmationTokenRepository {
	return &PostgresConfirmationTokenRepository{DB: db}
}

func (r *PostgresConfirmationTokenRepository) CreateToken(token *models.ConfirmationToken) error {
	query := `
		INSERT INTO confirmation_tokens (user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
	`
	token.CreatedAt = time.Now()
	_, err := r.DB.Exec(query, token.UserID, token.Token, token.ExpiresAt, token.CreatedAt)
	if err != nil {
		logger.Error("Error creating confirmation token: ", err)
	}
	return err
}

func (r *PostgresConfirmationTokenRepository) GetTokenByString(token string) (*models.ConfirmationToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at 
		FROM confirmation_tokens 
		WHERE token = $1
	`
	ct := &models.ConfirmationToken{}
	err := r.DB.QueryRow(query, token).Scan(
		&ct.ID,
		&ct.UserID,
		&ct.Token,
		&ct.ExpiresAt,
		&ct.CreatedAt,
	)
	if err != nil {
		logger.Error("Error fetching confirmation token: ", err)
		return nil, err
	}
	return ct, nil
}

func (r *PostgresConfirmationTokenRepository) DeleteToken(token string) error {
	query := `DELETE FROM confirmation_tokens WHERE token = $1`
	_, err := r.DB.Exec(query, token)
	if err != nil {
		logger.Error("Error deleting confirmation token: ", err)
	}
	return err
}
