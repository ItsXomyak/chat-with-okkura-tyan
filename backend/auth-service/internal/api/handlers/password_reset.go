package handlers

import (
	"encoding/json"
	"net/http"

	"authforge/internal/logger"
	"authforge/internal/services"
)

type PasswordResetHandler struct {
	AuthService services.AuthService
}

func NewPasswordResetHandler(authService services.AuthService) *PasswordResetHandler {
	return &PasswordResetHandler{
		AuthService: authService,
	}
}

type RequestResetRequest struct {
	Email string `json:"email"`
}

type RequestResetResponse struct {
	Message string `json:"message"`
}

func (h *PasswordResetHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	logger.Info("Password reset request received")
	var req RequestResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Invalid request payload: ", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		logger.Error("Email is required for password reset")
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	err := h.AuthService.RequestPasswordReset(req.Email)
	if err != nil {
		logger.Error("Request password reset failed: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Info("Password reset instructions sent for email: ", req.Email)
	response := RequestResetResponse{
		Message: "If this email is registered, password reset instructions have been sent.",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
}

func (h *PasswordResetHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	logger.Info("Reset password request received")
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Invalid request payload: ", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Token == "" || req.NewPassword == "" {
		logger.Error("Token and new password are required")
		http.Error(w, "Token and new password are required", http.StatusBadRequest)
		return
	}

	err := h.AuthService.ResetPassword(req.Token, req.NewPassword)
	if err != nil {
		logger.Error("Reset password failed: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Info("Password reset successfully")
	response := ResetPasswordResponse{
		Message: "Password has been reset successfully.",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
