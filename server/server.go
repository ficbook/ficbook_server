package server

import (
	"log"
	"net/http"
	"github.com/jinzhu/gorm"
	"github.com/gorilla/websocket"
	"github.com/gorilla/mux"
)

// Chat server.
type Server struct {
	messages    []*Message
	tempClients map[int]*Client
	clients   	map[string]*Client
	rooms		map[int]*Room
	addCh   	chan *Client
	delCh  		chan *Client
	sendQuery 	chan *Message
	doneCh	    chan bool
	errCh   	chan error
	db			*gorm.DB
}

// Create new chat server.
func NewServer(db *gorm.DB) *Server {
	var roomsSQL []*Room

	messages := []*Message{}
	tempClients := make(map[int]*Client)
	clients := make(map[string]*Client)
	rooms := make(map[int]*Room)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendQuery := make(chan *Message)
	doneCh := make(chan bool)
	errCh := make(chan error)

	db.Find(&roomsSQL)
	for _, room := range(roomsSQL) {
		rooms[room.ID] = NewRoom(room.ID, room.Name, room.CreatedAt)
	}

	return &Server{
		messages,
		tempClients,
		clients,
		rooms,
		addCh,
		delCh,
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

func (s *Server) SendQuery(i *Message) {
	s.sendQuery <- i
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

func (s *Server) sendToClient(msg *Message) {
	msg.Client.Write(msg)
}

func updrageSocket(w *http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("onConnected: %s", r.(error))
		}
	}()
	return upgrader.Upgrade(*w, r, nil)
}

// Listen and serve.
// It serves client connection and broadcast request.
func (s *Server) Listen() {
	log.Println("Listening server...")

/*	wsListen := func(w http.ResponseWriter, r *http.Request) {
		c, err := updrageSocket(&w, r)
		if err == nil {
			client := NewClient(c, s, s.db, RandStringRunes(8))
			s.Add(client)
			client.Listen()
		}
		
	}*/

	router := mux.NewRouter()
	router.HandleFunc("/users/sign_in", s.Users_SignIn).Methods("POST")
	router.HandleFunc("/rooms", s.Rooms_List).Methods("GET")
	router.HandleFunc("/rooms/{id}", s.Rooms_GetRoom).Methods("GET")
	router.HandleFunc("/rooms/new", s.Rooms_Create).Methods("POST")
	http.Handle("/", router)	

//	http.HandleFunc("/ws", wsListen)

	log.Println("Created handler")

	for {
		select {

		// Add new a client
		case c := <-s.addCh:
			s.tempClients[c.id] = c
			s.sendPastMessages(c)

		// del a client
		case c := <-s.delCh:
			delete(s.tempClients, c.id)

		case v := <-s.sendQuery:
			s.sendToClient(v)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}
