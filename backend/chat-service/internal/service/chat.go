package service

import (
	"strings"

	"chat-service/internal/db"
	"chat-service/internal/model"
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
  if err != nil {
    return nil, err
  }
  // Заглушка для AI
  if !isFromAI {
    aiContent := "Йо, я Okkura-тян! Классно болтать! 😎"
    if strings.Contains(content, "привет") {
      aiContent = "Привет-привет! Как дела, мой аниме-френд? 😺";
    }
    s.SendMessage(chatID, "okkura-tyan", aiContent, true)
  }
  return msg, nil
}

func (s *ChatService) GetMessages(chatID int) ([]model.Message, error) {
  var messages []model.Message
  err := s.db.Select(&messages, "SELECT * FROM messages WHERE chat_id = $1 ORDER BY created_at", chatID)
  return messages, err
}

func (s *ChatService) GetChats(userID string) ([]model.Chat, error) {
  var chats []model.Chat
  err := s.db.Select(&chats, "SELECT * FROM chats WHERE user_id = $1 ORDER BY created_at", userID)
  return chats, err
}

func (s *ChatService) UpdateMessage(id int, content string) (*model.Message, error) {
  msg := &model.Message{}
  err := s.db.QueryRowx("UPDATE messages SET content = $1 WHERE id = $2 RETURNING *", content, id).StructScan(msg)
  return msg, err
}

func (s *ChatService) DeleteMessage(id int) error {
  _, err := s.db.Exec("DELETE FROM messages WHERE id = $1", id)
  return err
}