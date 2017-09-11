package server

import (
	"time"
)

type Message struct {
	Client *Client
	Message *map[string]interface{}
}

type ChatMessage struct {
	ID uint32 `gorm:"primary_key; not null; AUTO_INCREMENT" json:"id"`
	UserID uint32 `json:"-" gorm:"not null"`
	UserName string `sql:"-" gorm:"not null"`
	Body string `json:"body" gorm:"not null"`
	RoomID int `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}

func (c *ChatMessage) TableName() string {
	return "messages"
}
