package model

import "time"

type Donation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"index;not null" json:"user_id"`
	Amount    float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Method    string    `gorm:"not null" json:"method"`
	Message   string    `json:"message,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (Donation) TableName() string {
	return "donates"
}
