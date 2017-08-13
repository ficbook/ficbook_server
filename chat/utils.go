package chat

import (
	"strings"
	"crypto/rand"
	"fmt"
	"time"
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

func (s *Server) GetSpecialRoomByName(name string) *Room {
	var index int
	for i, room := range(s.rooms) {
		if (*room).Name == name {
			index = i
			break
		}
	}
	return s.rooms[index]
}

func (s *Server) GetSpecialRoomByUUID(uuid string) *Room {
	var index int
	for i, room := range(s.rooms) {
		if strings.Contains((*room).UUID, uuid) {
			index = i
			break
		}
	}
	return s.rooms[index]
}

func RemoveAt(client *Client, room *Room) {
	index := -1
	for i, user := range(room.Users) {
		 if user.roomUUID == client.roomUUID {
			index = i
			break
		 }
	}
	copy(room.Users[index:], room.Users[index+1:])
	room.Users[len(room.Users)-1] = nil // or the zero value of T
	room.Users = room.Users[:len(room.Users)-1]
	fmt.Println(room.Users)
}

func (room *Room) RemoveByIndex(index int) {
	copy(room.Users[index:], room.Users[index+1:])
	room.Users[len(room.Users)-1] = nil // or the zero value of T
	room.Users = room.Users[:len(room.Users)-1]
	fmt.Println(room.Users)
}

func (s *Server) RefreshRoom() {
	rooms := []*Room{}
	var roomsSQL []*Room
	s.db.Table("chat_rooms").Find(&roomsSQL)
	for _, room := range roomsSQL {
		rooms = append(rooms, NewRoom(room.Id, room.Name, room.Topic, room.About, room.Type, room.UUID))
	}
	s.rooms = rooms
}

func (s *Server) UpdateOnline(updateTime int) {
	updateLastTime := time.Duration(updateTime) * time.Millisecond
	for {
		fmt.Println("[Updating count of users]")
		var count int
		for _, room := range(s.rooms) {
			count = 0
			for i, user := range(room.Users) {
				if user.ws.IsServerConn() {
					count++
				} else {
					room.RemoveByIndex(i)
				}
			}
			fmt.Printf(`Count of users in the room "%s" (before/now): `, room.Name)
			fmt.Printf("[%d/", room.LenUsers)
			fmt.Printf("%d]\n", count)
			room.LenUsers = count
		}
		fmt.Println("=====================")
		time.Sleep(updateLastTime)
	}
}