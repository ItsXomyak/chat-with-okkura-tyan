package cmd

import (
	"log"
	"net/http"

	"github.com/rs/cors"

	"auth-service/config"
	"auth-service/internal/api/handlers"
	"auth-service/internal/api/handlers/routes"
	"auth-service/internal/logger"
	"auth-service/internal/mailer"
	"auth-service/internal/repository"
	"auth-service/internal/services"
)

func Run() {
	logger.Init()

	cfg, err := config.LoadConfig(".")
	if err != nil {
		logger.Error("Error loading config: ", err)
		log.Fatalf("Error loading config: %v", err)
	}

	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		logger.Error("Error connecting to database: ", err)
		log.Fatalf("Error connecting to database: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewConfirmationTokenRepository(db)
	passwordResetTokenRepo := repository.NewPasswordResetTokenRepository(db)

	smtpMailer := mailer.NewSMTPMailer(cfg)

	authService := services.NewAuthService(userRepo, tokenRepo, passwordResetTokenRepo, cfg, smtpMailer)

	authHandler := handlers.NewAuthHandler(authService)
	confirmHandler := handlers.NewConfirmHandler(authService)
	passwordResetHandler := handlers.NewPasswordResetHandler(authService)

	routes.RegisterRoutes(authHandler, confirmHandler, passwordResetHandler)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(http.DefaultServeMux)

	logger.Info("Server starting on port ", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, corsHandler); err != nil {
		logger.Error("Server failed: ", err)
		log.Fatalf("Server failed: %v", err)
	}
}
