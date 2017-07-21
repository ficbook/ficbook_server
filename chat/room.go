package chat

import u "github.com/satori/go.uuid"
import "fmt"

type Room struct {
	Id int
	Name string
	Topic string
	UUID string	
}

//NewRoom returns the address of the room
func NewRoom(id int, name string, topic string, UUID string) *Room {
	return &Room{
		id,
		name,
		topic,
		UUID,
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
	}
}

