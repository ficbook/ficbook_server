package chat

import (
	"time"
)

type UserInfo struct {
	//id int `gorm:"primary_key"`
	Login []byte 
	Password []byte
	Power int
	DateReg time.Time
	DateVisit time.Time
}

func (u *UserInfo) TableName() string {
	return "users"
}