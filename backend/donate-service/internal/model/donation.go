package model

import "time"

type Donation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    string    `json:"user_id"`
	Amount    float64   `json:"amount"`
	Method    string    `json:"method"`
	Message   string    `json:"message,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
