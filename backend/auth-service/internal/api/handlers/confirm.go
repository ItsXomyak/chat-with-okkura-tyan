package handlers

import (
	"net/http"

	"auth-service/internal/logger"
	"auth-service/internal/services"
)

type ConfirmHandler struct {
	AuthService services.AuthService
}

func NewConfirmHandler(authService services.AuthService) *ConfirmHandler {
	return &ConfirmHandler{
		AuthService: authService,
	}
}

func (h *ConfirmHandler) ConfirmAccount(w http.ResponseWriter, r *http.Request) {
	logger.Info("Confirm account request received")
	token := r.URL.Query().Get("token")
	if token == "" {
		logger.Error("Token is missing in confirmation request")
		http.Error(w, "token is required", http.StatusBadRequest)
		return
	}

	err := h.AuthService.ConfirmAccount(token)
	if err != nil {
		logger.Error("Account confirmation failed: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Info("Account confirmed successfully")
	http.Redirect(w, r, "http://127.0.0.1:5502/static/email-confirmed.html", http.StatusSeeOther)
}
