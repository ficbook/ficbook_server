package chat

import (
	"time"
)

type Ban struct {
	LoginBanned []byte `json:"login_banned" gorm:"not null"`
	LoginBanning []byte `json:"login_baning" gorm:"not null"`
	Reason string `json:"reason" gorm:"not null"`
	TimeBan time.Time `json:"-" gorm:"not null"`
	TimeExpired time.Time `json:"time_expired" gorm:"not null"`
}

func (b *Ban) TableName() string {
	return "bans_list"
}

type AdminHistory struct {
	LoginAdmin []byte
	LoginUser []byte
	Action string
	Date time.Time
}