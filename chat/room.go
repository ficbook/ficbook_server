package chat

import (
	"fmt"
	"strings"
	u "github.com/satori/go.uuid"
)

type Room struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Topic string `json:"topic"`
	About string `json:"about"`
	UUID string `json:"uuid"`
	Users []*Client `json:"users"`
	LenUsers int `json:"count_users"`
}

//NewRoom returns the address of the room
func NewRoom(id int, name string, topic string, about, UUID string) *Room {
	return &Room{
		id,
		name,
		topic,
		about,
		UUID,
		[]*Client{},
		0,
	}
}

//CreateRoom returns the address of the created room
func CreateRoom(id int, name string, topic string, about string) *Room {
	uuid := u.NewV4()
	return &Room{
		id,
		name,
		topic,
		about,
		fmt.Sprint(uuid),
		[]*Client{},
		0,
	}
}

func GetSpecialRoomByName(rooms []*Room, name string) *Room {
	var index *int
	for i, room := range(rooms) {
		if strings.Contains(room.Name, name) {
			*index = i
			break
		}
	}
	return rooms[*index]
}
