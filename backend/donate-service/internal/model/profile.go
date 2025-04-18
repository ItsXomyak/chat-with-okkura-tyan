package model

import (
	"time"
)

type UserProfile struct {
	UserID    string    `gorm:"primaryKey" json:"user_id"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Role      string    `gorm:"default:user" json:"role"`
	Nickname  string    `json:"nickname"`
	FirstName string    `gorm:"column:first_name" json:"first_name"`
	LastName  string    `gorm:"column:last_name" json:"last_name"`
	AvatarURL string    `gorm:"column:avatar_url" json:"avatar_url"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (UserProfile) TableName() string {
	return "user_profiles"
}