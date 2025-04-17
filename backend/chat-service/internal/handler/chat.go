package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"chat-service/internal/service"
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
} // 123

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

func (h *ChatHandler) GetChats(c *gin.Context) {
  userID := c.Query("user_id")
  if userID == "" {
    c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
    return
  }
  chats, err := h.service.GetChats(userID)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }
  c.JSON(http.StatusOK, chats)
}

func (h *ChatHandler) UpdateMessage(c *gin.Context) {
  id, err := strconv.Atoi(c.Param("id"))
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
    return
  }
  var req struct {
    Content string `json:"content"`
  }
  if err := c.BindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }
  msg, err := h.service.UpdateMessage(id, req.Content)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }
  c.JSON(http.StatusOK, msg)
}

func (h *ChatHandler) DeleteMessage(c *gin.Context) {
  id, err := strconv.Atoi(c.Param("id"))
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
    return
  }
  if err := h.service.DeleteMessage(id); err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }
  c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}