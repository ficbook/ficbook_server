package chat

import (
	"log"
	"net/http"
	"time"
	"strings"
	"golang.org/x/net/websocket"
	"encoding/json"
	"github.com/jinzhu/gorm"
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
func NewServer(pattern string, db *gorm.DB) *Server {
	messages := []*Message{}
	clients := make(map[int]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan *Message)
	sendQuery := make(chan *InfoQuery)
	doneCh := make(chan bool)
	errCh := make(chan error)

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
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()

		client := NewClient(ws, s, s.db)
		s.Add(client)
		client.Listen()
	}
	http.Handle(s.pattern, websocket.Handler(onConnected))
	log.Println("Created handler")

	for {
		select {

		// Add new a client
		case c := <-s.addCh:
			s.clients[c.id] = c
			log.Println("Now", len(s.clients), "clients connected.")
			//s.sendPastMessages(c)

		// del a client
		case c := <-s.delCh:
			delete(s.clients, c.id)

		// broadcast message for all clients
		case msg := <-s.sendAllCh:
			//log.Println("Send all:", msg)
			s.messages = append(s.messages, msg)
			s.sendAll(msg)

		case v := <-s.sendQuery:
			ar := *v.ApiReturn
			client := *v.Client
			var m Message
			if ar.Interface == nil {
				interf := make(map[string]interface{})
				json.Unmarshal([]byte(ar.Text), &interf)
				m = Message{ar.Type, interf}
			} else {
				m = Message{ar.Type, *(ar.Interface)}
			}
			log.Print(m)
			if ar.ReturnVariable != nil {
				if ar.ReturnVariable.code == 7777 {
					s.sendAll(&m)
				} else {
					s.sendToClient(v.Client, &m)
				}
			} else {
				s.sendToClient(v.Client, &m)
			}
			if strings.Contains(ar.Type, "AUTH_ERROR") {
				time.Sleep(1)
				client.ws.Close()
				delete(s.clients, client.id)
			}

	//	case err := <-s.errCh:
	//		continue

		case <-s.doneCh:
			return
		}
	}
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