package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"otaku-verse/chat-service/internal/model"
	"otaku-verse/chat-service/internal/service"
)

var upgrader = websocket.Upgrader{
  CheckOrigin: func(r *http.Request) bool {
    return r.Header.Get("Origin") == "http://127.0.0.1:5501" // Live Server
  },
}

type WebSocketHandler struct {
  service *service.ChatService
  clients map[int]map[*websocket.Conn]bool
}

func NewWebSocketHandler(service *service.ChatService) *WebSocketHandler {
  return &WebSocketHandler{
    service: service,
    clients: make(map[int]map[*websocket.Conn]bool),
  }
}

func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
  chatID := c.Param("id")
  conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
  if err != nil {
    log.Println("WebSocket upgrade:", err)
    return
  }

  chatIDInt, _ := strconv.Atoi(chatID)
  if h.clients[chatIDInt] == nil {
    h.clients[chatIDInt] = make(map[*websocket.Conn]bool)
  }
  h.clients[chatIDInt][conn] = true

  defer func() {
    delete(h.clients[chatIDInt], conn)
    conn.Close()
  }()

  for {
    var msg model.Message
    if err := conn.ReadJSON(&msg); err != nil {
      log.Println("WebSocket read:", err)
      break
    }
    msg.ChatID = chatIDInt
    savedMsg, err := h.service.SendMessage(msg.ChatID, msg.UserID, msg.Content, msg.IsFromAI)
    if err != nil {
      log.Println("Save message:", err)
      continue
    }
    for client := range h.clients[chatIDInt] {
      client.WriteJSON(savedMsg)
    }
  }
}