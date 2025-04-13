package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"otaku-verse/chat-service/internal/db"
	"otaku-verse/chat-service/internal/handler"
	"otaku-verse/chat-service/internal/service"
)

func main() {
  // Миграции
  m, err := migrate.New("file://migrations", "postgres://user:password@localhost:5432/chatdb?sslmode=disable")
  if err != nil {
    log.Fatal("migrate init:", err)
  }
  if err := m.Up(); err != nil && err != migrate.ErrNoChange {
    log.Fatal("migrate up:", err)
  }

  // База
  dbConn, err := db.New("postgres://user:password@localhost:5432/chatdb?sslmode=disable")
  if err != nil {
    log.Fatal("db init:", err)
  }
  defer dbConn.Close()

  // Сервис и хендлеры
  chatService := service.NewChatService(dbConn)
  chatHandler := handler.NewChatHandler(chatService)

  // Роутер
  r := gin.Default()
  r.POST("/chats", chatHandler.CreateChat)
  r.POST("/chats/:id/messages", chatHandler.SendMessage)
  r.GET("/chats/:id/messages", chatHandler.GetMessages)

  // Запуск
  if err := r.Run(":8080"); err != nil {
    log.Fatal("server run:", err)
  }
}