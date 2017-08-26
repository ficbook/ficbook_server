package chat

import (
	"time"
)

type UserInfo struct {
	ID int `gorm:"primary_key"`
	Login []byte `gorm:"type:varbinary(255);not null"`
	Password []byte `gorm:"type:varbinary(255);not null"`
	Power int `gorm:"not null"`
	DateReg time.Time `gorm:"not null"`
	DateVisit time.Time `gorm:"not null"`
}

func (u *UserInfo) TableName() string {
	return "users"
}