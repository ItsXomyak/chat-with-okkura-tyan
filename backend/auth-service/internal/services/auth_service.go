package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"authforge/config"
	"authforge/internal/logger"
	"authforge/internal/mailer"
	"authforge/internal/models"
	"authforge/internal/repository"
)

type AuthService interface {
	RegisterUser(user *models.User, password string) error
	Login(email, password string) (*TokenPair, error)
	ConfirmAccount(tokenString string) error
	RequestPasswordReset(email string) error
	ResetPassword(token, newPassword string) error
	ValidateToken(tokenString string) (*models.CustomClaims, error)
}

type authService struct {
	userRepo               repository.UserRepository
	tokenRepo              repository.ConfirmationTokenRepository
	passwordResetTokenRepo repository.PasswordResetTokenRepository
	cfg                    *config.Config
	mailer                 mailer.Mailer
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func NewAuthService(
	userRepo repository.UserRepository,
	tokenRepo repository.ConfirmationTokenRepository,
	passwordResetTokenRepo repository.PasswordResetTokenRepository,
	cfg *config.Config,
	m mailer.Mailer,
) AuthService {
	logger.Info("Initializing AuthService")
	return &authService{
		userRepo:               userRepo,
		tokenRepo:              tokenRepo,
		passwordResetTokenRepo: passwordResetTokenRepo,
		cfg:                    cfg,
		mailer:                 m,
	}
}

func (s *authService) RegisterUser(user *models.User, password string) error {
	existingUser, err := s.userRepo.GetUserByEmail(user.Email)
	if err == nil && existingUser != nil {
		logger.Error("User already exists: ", user.Email)
		return errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Error hashing password for ", user.Email, ": ", err)
		return err
	}

	user.PasswordHash = string(hashedPassword)
	user.IsActive = false

	if user.Role == "" {
		user.Role = "user"
	}

	if user.Role != models.RoleUser && user.Role != models.RoleAdmin {
		logger.Error("Invalid role specified for user: ", user.Email)
		return errors.New("invalid role")
	}

	user.ID = uuid.New()

	if err := s.userRepo.CreateUser(user); err != nil {
		logger.Error("Error creating user ", user.Email, ": ", err)
		return err
	}

	confirmationToken, err := generateRandomToken(32)
	if err != nil {
		logger.Error("Error generating confirmation token: ", err)
		return err
	}

	token := &models.ConfirmationToken{
		UserID:    user.ID,
		Token:     confirmationToken,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.tokenRepo.CreateToken(token); err != nil {
		logger.Error("Error saving confirmation token for user ", user.Email, ": ", err)
		return err
	}

	if err := s.mailer.SendConfirmationEmail(user.Email, confirmationToken); err != nil {
		logger.Error("Error sending confirmation email to ", user.Email, ": ", err)
		return err
	}

	return nil
}

func (s *authService) Login(email, password string) (*TokenPair, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		logger.Error("Login failed, user not found: ", email)
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		logger.Error("Login failed, account not activated: ", email)
		return nil, errors.New("account not activated")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		logger.Error("Login failed, invalid credentials for: ", email)
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := s.generateJWTToken(user, s.cfg.JWTExpiry)
	if err != nil {
		logger.Error("Error generating access token for ", email, ": ", err)
		return nil, err
	}

	refreshToken, err := s.generateJWTToken(user, s.cfg.RefreshExpiry)
	if err != nil {
		logger.Error("Error generating refresh token for ", email, ": ", err)
		return nil, err
	}
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) generateJWTToken(user *models.User, expiry time.Duration) (string, error) {
	claims := &models.CustomClaims{
		UserID: user.ID.String(),
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func generateRandomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		logger.Error("Error generating random bytes: ", err)
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *authService) ConfirmAccount(tokenString string) error {
	confirmationToken, err := s.tokenRepo.GetTokenByString(tokenString)
	if err != nil {
		logger.Error("Invalid confirmation token: ", err)
		return errors.New("invalid token")
	}

	if time.Now().After(confirmationToken.ExpiresAt) {
		logger.Error("Confirmation token expired for user ", confirmationToken.UserID)
		return errors.New("token expired")
	}

	user, err := s.userRepo.GetUserByID(confirmationToken.UserID)
	if err != nil {
		logger.Error("Error retrieving user for confirmation: ", err)
		return err
	}

	user.IsActive = true
	if err := s.userRepo.UpdateUser(user); err != nil {
		logger.Error("Error updating user status for ", user.Email, ": ", err)
		return err
	}

	if err := s.tokenRepo.DeleteToken(tokenString); err != nil {
		logger.Error("Error deleting confirmation token: ", err)
	}

	return nil
}

func (s *authService) RequestPasswordReset(email string) error {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		logger.Error("User not found for password reset: ", email)
		return errors.New("user not found")
	}

	resetToken, err := generateRandomToken(32)
	if err != nil {
		logger.Error("Error generating password reset token for ", email, ": ", err)
		return err
	}

	tokenModel := &models.PasswordResetToken{
		UserID:    user.ID,
		Token:     resetToken,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
	}

	if err := s.passwordResetTokenRepo.CreateToken(tokenModel); err != nil {
		logger.Error("Error saving password reset token for ", email, ": ", err)
		return err
	}

	if err := s.mailer.SendPasswordResetEmail(user.Email, resetToken); err != nil {
		logger.Error("Error sending password reset email to ", email, ": ", err)
		return err
	}

	return nil
}

func (s *authService) ResetPassword(tokenStr, newPassword string) error {
	tokenModel, err := s.passwordResetTokenRepo.GetToken(tokenStr)
	if err != nil {
		logger.Error("Invalid password reset token: ", err)
		return errors.New("invalid token")
	}

	if tokenModel.Used {
		logger.Error("Password reset token already used for user ", tokenModel.UserID)
		return errors.New("token already used")
	}

	if time.Now().After(tokenModel.ExpiresAt) {
		logger.Error("Password reset token expired for user ", tokenModel.UserID)
		return errors.New("token expired")
	}

	user, err := s.userRepo.GetUserByID(tokenModel.UserID)
	if err != nil {
		logger.Error("Error retrieving user for password reset: ", err)
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Error hashing new password for user ", user.Email, ": ", err)
		return err
	}
	user.PasswordHash = string(hashedPassword)

	if err := s.userRepo.UpdateUser(user); err != nil {
		logger.Error("Error updating password for user ", user.Email, ": ", err)
		return err
	}

	if err := s.passwordResetTokenRepo.MarkTokenUsed(tokenStr); err != nil {
		logger.Error("Error marking password reset token as used: ", err)
		return err
	}

	return nil
}

func (s *authService) ValidateToken(tokenString string) (*models.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
