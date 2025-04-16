package handlers

import (
	"encoding/json"
	"net/http"

	"authforge/internal/logger"
	"authforge/internal/models"
	"authforge/internal/services"
)

type AuthHandler struct {
	AuthService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type ResponseMessage struct {
	Message string `json:"message"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	logger.Info("Registration request received")
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Invalid request payload: ", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" {
		logger.Error("Email and password are required")
		http.Error(w, "Email and password required", http.StatusBadRequest)
		return
	}

	user := &models.User{
		Email: req.Email,
		Role:  models.UserRole(req.Role),
	}

	if err := h.AuthService.RegisterUser(user, req.Password); err != nil {
		logger.Error("Registration failed: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("User registered successfully: ", req.Email)
	resp := ResponseMessage{Message: "Registration successful. Please check your email to activate your account."}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	logger.Info("Login request received")
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Invalid request payload: ", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" {
		logger.Error("Email and password are required for login")
		http.Error(w, "Email and password required", http.StatusBadRequest)
		return
	}

	tokens, err := h.AuthService.Login(req.Email, req.Password)
	if err != nil {
		logger.Error("Login failed for ", req.Email, ": ", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	resp := LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	logger.Info("User logged in successfully: ", req.Email)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
