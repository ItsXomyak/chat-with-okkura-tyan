package main

import (
	"log"

	"otaku-verse/chat-service/internal/db"
	"otaku-verse/chat-service/internal/handler"
	"otaku-verse/chat-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
  m, err := migrate.New("file://migrations", "postgres://user:password@localhost:5432/chatdb?sslmode=disable")
  if err != nil {
    log.Fatal("migrate init:", err)
  }
  if err := m.Up(); err != nil && err != migrate.ErrNoChange {
    log.Fatal("migrate up:", err)
  }

  dbConn, err := db.New("postgres://user:password@localhost:5432/chatdb?sslmode=disable")
  if err != nil {
    log.Fatal("db init:", err)
  }
  defer dbConn.Close()

  chatService := service.NewChatService(dbConn)
  chatHandler := handler.NewChatHandler(chatService)
  wsHandler := handler.NewWebSocketHandler(chatService)

  r := gin.Default()

  r.Use(func(c *gin.Context) {
    c.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5501")
    c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
    c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    if c.Request.Method == "OPTIONS" {
      c.AbortWithStatus(204)
      return
    }
    c.Next()
  })

  r.GET("/chats", chatHandler.GetChats)
  r.POST("/chats", chatHandler.CreateChat)
  r.POST("/chats/:id/messages", chatHandler.SendMessage)
  r.GET("/chats/:id/messages", chatHandler.GetMessages)
  r.PATCH("/messages/:id", chatHandler.UpdateMessage)
  r.DELETE("/messages/:id", chatHandler.DeleteMessage)
  r.GET("/ws/chats/:id", wsHandler.HandleWebSocket)

  if err := r.Run(":8080"); err != nil {
    log.Fatal("server run:", err)
  }
}