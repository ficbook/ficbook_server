package chat

import u "github.com/satori/go.uuid"
import "fmt"

type Room struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Topic string `json:"topic"`
	UUID string `json:"uuid"`
	Users []*Client `json:"users"`
}

//NewRoom returns the address of the room
func NewRoom(id int, name string, topic string, UUID string) *Room {
	return &Room{
		id,
		name,
		topic,
		UUID,
		[]*Client{},
	}
}

//CreateRoom returns the address of the created room
func CreateRoom(id int, name string, topic string) *Room {
	uuid := u.NewV4()
	return &Room{
		id,
		name,
		topic,
		fmt.Sprint(uuid),
		[]*Client{},
	}
}

