package chat

import (
	"time"
)

type User struct {
	login []byte
	password []byte
	power int
	DateReg time.Time
	DateVisit time.Time
}

func NewUserString(login string, password string, power int, datereg time.Time, datevisit time.Time) *User{
	return &User{
		[]byte(login),
		[]byte(password),
		power,
		datereg,
		datevisit,	
	}
}

func NewUserByte(login []byte, password []byte, power int, datereg time.Time, datevisit time.Time) *User{
	return &User{
		[]byte(login),
		[]byte(password),
		power,
		datereg,
		datevisit,	
	}
}