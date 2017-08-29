package chat

import (
	"strings"
	"crypto/rand"
	"fmt"
	"time"
	"sort"
	u "github.com/satori/go.uuid"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
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

func (room *Room) RemoveAt(clientID int) {
	delete(room.Users, clientID)
}

func (s *Server) RefreshRoom() {
	rooms := make(map[int]*Room)
	var roomsSQL []*Room
	s.db.Find(&roomsSQL)
	for _, room := range roomsSQL {
		rooms[room.ID] = NewRoom(room.ID, room.Name, room.Topic, room.About, room.Type, room.UUID)
	}
	s.rooms = rooms
}

func GetLoginUsers(users map[int]*Client) *[]string {
	var logins []string
	for _, u := range(users) {
		logins = append(logins, string(u.userInfo.Login))
	}
	if len(logins) == 0 {
		logins = []string{}
	}
	return &logins
}

func (s *Server) UpdateOnline(updateTime int) {
	updateLastTime := time.Duration(updateTime) * time.Millisecond
	for {
		var count int
		for _, room := range(s.rooms) {
			count = 0
			for _, user := range(room.Users) {
				checkMap := NewMap()
				(*checkMap)["type"] = "check"
				err := user.ws.WriteJSON(checkMap)
				if err != nil {
					fmt.Print("ERROR: ")
					fmt.Println(err)
					user.ws.Close()
					room.RemoveAt(user.id)
				} else {
					user.antiflood = 0
					count++
				}
			}
			room.LenUsers = count
		}
		time.Sleep(updateLastTime)	
	}
}

func (s *Server) SearchUser(userLogin string) (*Client, bool) {
	isSearch := false
	client := Client{}
	for _, c := range(s.clients) {
		if c.StringLogin() == userLogin {
			client = *c
			isSearch = true
		}
	}
	return &client, isSearch
}

func (s *Server) UpdateListRooms() {
	var(
		roomInts []int
		roomsList []*Room
	)

	for k := range(s.rooms) {
		roomInts = append(roomInts, k)
	}

	sort.Ints(roomInts)
	
	for _, k := range(roomInts) {
		roomsList = append(roomsList, s.rooms[k])
	}

	s.roomsList = &roomsList
}