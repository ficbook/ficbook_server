package server

import "time"

type User struct {
	ID uint32 `gorm:"primary_key; not null; AUTO_INCREMENT" json:"id"`
	Name string `json:"name" gorm:"not null"`
	Token string `json:"token" gorm:"not null"`
	CreatedAt time.Time `json:"-" gorm:"not null"`
	LastLogin time.Time `json:"last_login"`
}