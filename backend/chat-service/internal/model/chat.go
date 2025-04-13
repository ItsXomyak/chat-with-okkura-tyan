package model

import "time"

type Chat struct {
  ID        int       `json:"id" db:"id"`
  UserID    string    `json:"user_id" db:"user_id"`
  Name      string    `json:"name" db:"name"`
  CreatedAt time.Time `json:"created_at" db:"created_at"`
}