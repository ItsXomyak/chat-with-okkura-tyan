package handlers

import (
	"encoding/json"
	"net/http"

	"authforge/internal/logger"
	"authforge/internal/services"
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
	response := map[string]string{"message": "Account activated successfully."}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
