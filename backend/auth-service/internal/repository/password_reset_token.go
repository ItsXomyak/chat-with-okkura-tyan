package repository

import (
	"database/sql"
	"time"

	"authforge/internal/logger"
	"authforge/internal/models"
)

type PasswordResetTokenRepository interface {
	CreateToken(token *models.PasswordResetToken) error
	GetToken(token string) (*models.PasswordResetToken, error)
	MarkTokenUsed(token string) error
}

type PostgresPasswordResetTokenRepository struct {
	DB *sql.DB
}

func NewPasswordResetTokenRepository(db *sql.DB) PasswordResetTokenRepository {
	return &PostgresPasswordResetTokenRepository{DB: db}
}

func (r *PostgresPasswordResetTokenRepository) CreateToken(token *models.PasswordResetToken) error {
	query := `
		INSERT INTO password_reset_tokens (user_id, token, expires_at, created_at, used)
		VALUES ($1, $2, $3, $4, $5)
	`
	token.CreatedAt = time.Now()
	token.Used = false
	_, err := r.DB.Exec(query, token.UserID, token.Token, token.ExpiresAt, token.CreatedAt, token.Used)
	if err != nil {
		logger.Error("Error creating password reset token for user ", token.UserID, ": ", err)
	}
	return err
}

func (r *PostgresPasswordResetTokenRepository) GetToken(tokenStr string) (*models.PasswordResetToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at, used
		FROM password_reset_tokens
		WHERE token = $1
	`
	prt := &models.PasswordResetToken{}
	err := r.DB.QueryRow(query, tokenStr).Scan(
		&prt.ID,
		&prt.UserID,
		&prt.Token,
		&prt.ExpiresAt,
		&prt.CreatedAt,
		&prt.Used,
	)
	if err != nil {
		logger.Error("Error fetching password reset token: ", err)
		return nil, err
	}
	return prt, nil
}

func (r *PostgresPasswordResetTokenRepository) MarkTokenUsed(tokenStr string) error {
	query := `UPDATE password_reset_tokens SET used = true WHERE token = $1`
	_, err := r.DB.Exec(query, tokenStr)
	if err != nil {
		logger.Error("Error marking password reset token as used: ", err)
	}
	return err
}
