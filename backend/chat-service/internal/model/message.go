package model

import "time"

type Message struct {
  ID        int       `json:"id" db:"id"`
  ChatID    int       `json:"chat_id" db:"chat_id"`
  UserID    string    `json:"user_id" db:"user_id"`
  Content   string    `json:"content" db:"content"`
  IsFromAI  bool      `json:"is_from_ai" db:"is_from_ai"`
  CreatedAt time.Time `json:"created_at" db:"created_at"`
}