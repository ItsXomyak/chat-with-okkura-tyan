package routes

import (
	"authforge/internal/api/handlers"
	"net/http"
)

func RegisterRoutes(
	authHandler *handlers.AuthHandler,
	confirmHandler *handlers.ConfirmHandler,
	passwordResetHandler *handlers.PasswordResetHandler,
) {
	http.HandleFunc("/api/v1/auth/register", authHandler.Register)
	http.HandleFunc("/api/v1/auth/login", authHandler.Login)
	http.HandleFunc("/api/v1/auth/confirm", confirmHandler.ConfirmAccount)
	http.HandleFunc("/api/v1/auth/password-reset-request", passwordResetHandler.RequestPasswordReset)
	http.HandleFunc("/api/v1/auth/password-reset-confirm", passwordResetHandler.ResetPassword)
	http.HandleFunc("/api/v1/auth/validate", authHandler.ValidateToken)
}
