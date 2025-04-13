package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"otaku-verse/chat-service/internal/service"
)

type ChatHandler struct {
  service *service.ChatService
}

func NewChatHandler(service *service.ChatService) *ChatHandler {
  return &ChatHandler{service}
}

func (h *ChatHandler) CreateChat(c *gin.Context) {
  var req struct {
    UserID string `json:"user_id"`
    Name   string `json:"name"`
  }
  if err := c.BindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }
  chat, err := h.service.CreateChat(req.UserID, req.Name)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }
  c.JSON(http.StatusOK, chat)
}

func (h *ChatHandler) SendMessage(c *gin.Context) {
  chatID, err := strconv.Atoi(c.Param("id"))
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
    return
  }
  var req struct {
    UserID   string `json:"user_id"`
    Content  string `json:"content"`
    IsFromAI bool   `json:"is_from_ai"`
  }
  if err := c.BindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }
  msg, err := h.service.SendMessage(chatID, req.UserID, req.Content, req.IsFromAI)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }
  c.JSON(http.StatusOK, msg)
}

func (h *ChatHandler) GetMessages(c *gin.Context) {
  chatID, err := strconv.Atoi(c.Param("id"))
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
    return
  }
  messages, err := h.service.GetMessages(chatID)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }
  c.JSON(http.StatusOK, messages)
}