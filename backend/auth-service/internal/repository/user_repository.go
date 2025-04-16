package repository

import (
	"database/sql"
	"errors"
	"time"

	"authforge/internal/logger"
	"authforge/internal/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
	UpdateUser(user *models.User) error
}

type PostgresUserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &PostgresUserRepository{DB: db}
}

func (r *PostgresUserRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (
			id, email, password_hash, is_active, role,
			created_at, updated_at, failed_login_attempts, last_failed_login
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	now := time.Now()
	user.ID = uuid.New() // ✅ явно создаём UUID
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := r.DB.Exec(query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.IsActive,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
		user.FailedLoginAttempts,
		user.LastFailedLogin,
	)

	if err != nil {
		logger.Error("Error creating user with email ", user.Email, ": ", err)
	}
	return err
}

func (r *PostgresUserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, is_active, role, created_at, updated_at, failed_login_attempts, last_failed_login
		FROM users WHERE email = $1`
	user := &models.User{}
	err := r.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.IsActive,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.FailedLoginAttempts,
		&user.LastFailedLogin,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("User not found with email ", email)
			return nil, errors.New("user not found")
		}
		logger.Error("Error fetching user by email ", email, ": ", err)
		return nil, err
	}
	return user, nil
}

func (r *PostgresUserRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, is_active, role, created_at, updated_at, failed_login_attempts, last_failed_login
		FROM users WHERE id = $1`
	user := &models.User{}
	err := r.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.IsActive,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.FailedLoginAttempts,
		&user.LastFailedLogin,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("User not found with ID ", id)
			return nil, errors.New("user not found")
		}
		logger.Error("Error fetching user by ID ", id, ": ", err)
		return nil, err
	}
	return user, nil
}

func (r *PostgresUserRepository) UpdateUser(user *models.User) error {
	query := `
		UPDATE users 
		SET email = $1, password_hash = $2, is_active = $3, role = $4, updated_at = $5, failed_login_attempts = $6, last_failed_login = $7
		WHERE id = $8`
	user.UpdatedAt = time.Now()
	_, err := r.DB.Exec(query,
		user.Email,
		user.PasswordHash,
		user.IsActive,
		user.Role,
		user.UpdatedAt,
		user.FailedLoginAttempts,
		user.LastFailedLogin,
		user.ID,
	)
	if err != nil {
		logger.Error("Error updating user with ID ", user.ID, ": ", err)
	}
	return err
}
