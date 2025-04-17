package model

type UserProfile struct {
	UserID    string `gorm:"primaryKey"`
	Nickname  string
	FirstName string
	LastName  string
	AvatarURL string
}