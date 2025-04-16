package handlers

import (
	"authforge/internal/logger"
	"authforge/internal/services"
	"encoding/json"
	"net/http"
	"strings"
)

type ValidateHandler struct {
	AuthService services.AuthService
}

func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	logger.Info("Token validation request received")
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		http.Error(w, "invalid token format", http.StatusUnauthorized)
		return
	}
	tokenStr := strings.TrimSpace(authHeader[len(prefix):])

	claims, err := h.AuthService.ValidateToken(tokenStr)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	response := map[string]interface{}{
		"user_id":   claims.UserID,
		"role":      claims.Role,
		"expiresAt": claims.ExpiresAt.Time,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
