package chat

import (
	"log"
	"net/http"
	//"container/list"
	"fmt"
	"golang.org/x/net/websocket"
)

// Chat server.
type Server struct {
	pattern   string
	messages  []*Message
	clients   map[int]*Client
	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan *Message
	sendQuery chan *InfoQuery
	doneCh    chan bool
	errCh     chan error
}

// Create new chat server.
func NewServer(pattern string) *Server {
	messages := []*Message{}
	clients := make(map[int]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan *Message)
	sendQuery := make(chan *InfoQuery)
	doneCh := make(chan bool)
	errCh := make(chan error)

	return &Server{
		pattern,
		messages,
		clients,
		addCh,
		delCh,
		sendAllCh,
		sendQuery,
		doneCh,
		errCh,
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

		client := NewClient(ws, s)
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
			s.sendPastMessages(c)
		//	s.sendToClient(c, &Message{""})
			id := fmt.Sprintf("%d", (*c).id)
			s.sendAll(&Message{Login: "[Сервер]", Text: "Пользователь присоединился к чату. ID: " + id})

		// del a client
		case c := <-s.delCh:
			id := fmt.Sprintf("%d", (*c).id)
			s.sendAll(&Message{Login: "[Сервер]", Text: "Пользователь вышел из чата. ID: " + id})
			delete(s.clients, c.id)

		// broadcast message for all clients
		case msg := <-s.sendAllCh:
			//log.Println("Send all:", msg)
			s.messages = append(s.messages, msg)
			s.sendAll(msg)

		case v := <-s.sendQuery:
			ar := *v.ApiReturn
			m := Message{ar.Type, ar.Text}
			s.sendToClient(v.Client, &m)

	//	case err := <-s.errCh:
	//		continue

		case <-s.doneCh:
			return
		}
	}
}
