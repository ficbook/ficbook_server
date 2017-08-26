package chat

import "time"

type Message struct {
	Type string `json:"type"`
	Text map[string]interface{}
}

func (self *Message) String() string {
	return "Send Message [" + self.Type + "]"
}

type ChatMessageSQL struct {
	ID int `gorm:"primary_key"`
	Login string `gorm:"unique_users_login; not null"`
	Message string
	Timestamp time.Time `sql:"timestamp"`
	RoomUUID string `gorm:"not null"`
}

func (c *ChatMessageSQL) TableName() string {
	return "chat_message_all"
}

type ChatMessageJSON struct {
	Login string `json:"login"`
	Message string `json:"message"`
	Timestamp int64 `json:"timestamp"`
}

func NewChatMessageJSON(login string, message string, timestamp time.Time) ChatMessageJSON {
	return ChatMessageJSON{
		login,
		message,
		timestamp.UnixNano() / 1000000,
	}
}