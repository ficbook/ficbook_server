package chat

import (
	"log"
	"net/http"
	"time"
	"strings"
	"github.com/jinzhu/gorm"
	"github.com/gorilla/websocket"
)

// Chat server.
type Server struct {
	pattern   string
	messages  []*Message
	rooms	  []*Room
	clients   map[int]*Client
	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan *Message
	sendQuery chan *InfoQuery
	doneCh    chan bool
	errCh     chan error
	db	*gorm.DB
}

// Create new chat server.
func NewServer(pattern string, db *gorm.DB, isRebuild *bool, createRoom *string) *Server {
	messages := []*Message{}
	clients := make(map[int]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan *Message)
	sendQuery := make(chan *InfoQuery)
	doneCh := make(chan bool)
	errCh := make(chan error)

	if *isRebuild {
		db.AutoMigrate(&UserInfo{},&ChatMessageSQL{})
		db.Table("chat_rooms").AutoMigrate(&Room{})
	}

	if len(*createRoom) > 0 {
		db.Table("chat_rooms").Create(CreateRoom(0, *createRoom, "", *createRoom, "public"))
	}

	rooms := []*Room{}
	var roomsSQL []*Room
	db.Table("chat_rooms").Find(&roomsSQL)
	for _, room := range roomsSQL {
		rooms = append(rooms, NewRoom(room.Id, room.Name, room.Topic, room.About, room.Type, room.UUID))
	}

	return &Server{
		pattern,
		messages,
		rooms,
		clients,
		addCh,
		delCh,
		sendAllCh,
		sendQuery,
		doneCh,
		errCh,
		db,
	}
}

var upgrader = websocket.Upgrader{}

func (s *Server) Add(c *Client) {
	s.addCh <- c
}

func (s *Server) Del(c *Client) {
	s.delCh <- c
}

func (s *Server) SendQuery(i *InfoQuery) {
	s.sendQuery <- i
}

func (s *Server) SendAll(msg *Message) {
	s.sendAllCh <- msg
}

func (s *Server) Done() {
	s.doneCh <- true
}

func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) sendPastMessages(c *Client) {
	for _, msg := range s.messages {
		c.Write(msg)
	}
}

func (s *Server) sendAll(msg *Message) {
	for _, c := range s.clients {
		c.Write(msg)
	}
}

func (s *Server) sendToClient(client *Client, msg *Message) {
	client.Write(msg)
}
// Listen and serve.
// It serves client connection and broadcast request.
func (s *Server) Listen() {

	log.Println("Listening server...")

	// websocket handler
	onConnected := func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatalf("onConnected", err)
			return
		}

		client := NewClient(c, s, s.db)
		s.Add(client)
		client.Listen()
	}
	http.HandleFunc(s.pattern, onConnected)
	log.Println("Created handler")

	for {
		select {

		// Add new a client
		case c := <-s.addCh:
			log.Println("Added new client")
			s.clients[c.id] = c
			log.Println("Now", len(s.clients), "clients connected.")
			s.sendPastMessages(c)

		// del a client
		case c := <-s.delCh:
			log.Println("Delete client")
			delete(s.clients, c.id)

		// broadcast message for all clients
		case msg := <-s.sendAllCh:
			log.Println("Send all:", msg)
			s.messages = append(s.messages, msg)
			s.sendAll(msg)

		case v := <-s.sendQuery:
			ar := *v.ApiReturn
			client := v.Client
			if !client.isAuth && ar.Type != "AUTH_OK" {
				client.ws.Close()
				delete(s.clients, client.id)
			}
			m := Message{ar.Type, *ar.Interface}
			if ar.ReturnVariable != nil {
				if ar.ReturnVariable.ReturnRoom != nil {
					if ar.ReturnVariable.int == 35 {
						for _, user := range(ar.ReturnVariable.ReturnRoom.Users) {
							if user.roomUUID == ar.ReturnVariable.ReturnRoom.Name {
								s.sendToClient(user, &m)
							}
						}
					}
				}
			} else {
				s.sendToClient(client, &m)
			}
			if strings.Contains(ar.Type, "AUTH_ERROR") {
				time.Sleep(1)
				client.ws.Close()
				delete(s.clients, client.id)
			}

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}

func (s Server) UpdateOnlineRooms(updateTime int) {
	s.UpdateOnline(updateTime)
}