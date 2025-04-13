package service

import (
	"otaku-verse/chat-service/internal/db"
	"otaku-verse/chat-service/internal/model"
)

type ChatService struct {
  db *db.DB
}

func NewChatService(db *db.DB) *ChatService {
  return &ChatService{db}
}

func (s *ChatService) CreateChat(userID, name string) (*model.Chat, error) {
  chat := &model.Chat{}
  err := s.db.QueryRowx("INSERT INTO chats (user_id, name) VALUES ($1, $2) RETURNING *", userID, name).StructScan(chat)
  return chat, err
}

func (s *ChatService) SendMessage(chatID int, userID, content string, isFromAI bool) (*model.Message, error) {
  msg := &model.Message{}
  err := s.db.QueryRowx("INSERT INTO messages (chat_id, user_id, content, is_from_ai) VALUES ($1, $2, $3, $4) RETURNING *", chatID, userID, content, isFromAI).StructScan(msg)
  return msg, err
}

func (s *ChatService) GetMessages(chatID int) ([]model.Message, error) {
  var messages []model.Message
  err := s.db.Select(&messages, "SELECT * FROM messages WHERE chat_id = $1 ORDER BY created_at", chatID)
  return messages, err
}