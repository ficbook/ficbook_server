package chat

import (
	"crypto/rand"
	"fmt"
	"strings"
	u "github.com/satori/go.uuid"
)

func GenerateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func NewUUID() string {
	return fmt.Sprint(u.NewV4())
}

func GetSpecialRoomByName(rooms []*Room, name string) *Room {
	var index int
	for i, room := range(rooms) {
		if strings.Contains(room.Name, name) {
			index = i
			break
		}
	}
	return rooms[index]
}
