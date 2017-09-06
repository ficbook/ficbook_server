package server

import (
	"time"
)

type Room struct {
	ID int `gorm:"primary_key; not null; AUTO_INCREMENT" json:"id"`
	Name string `gorm:"not null" json:"name"`
	CreatedAt time.Time `gorm:"not null" json:"-"`
	Members int `json:"members" sql:"-"`
}

func NewRoom(id int, name string, Time time.Time) *Room {
	return &Room{
		id,
		name,
		Time,
		0,
	}
}