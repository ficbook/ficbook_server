package chat

import (
	"fmt"
	"io"
	"log"
	"github.com/jinzhu/gorm"
	"github.com/gorilla/websocket"
)

const channelBufSize = 100

var maxId int = 0

// Chat client.
type Client struct {
	id     int
	ws     *websocket.Conn
	server *Server
	ch     chan *Message
	doneCh chan bool
	userInfo *UserInfo
	login string
	roomUUID string
	isAuth bool
	db *gorm.DB
}

// Create new chat client.
func NewClient(ws *websocket.Conn, server *Server, db *gorm.DB) *Client {

	if ws == nil {
		panic("ws cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}

	maxId++
	ch := make(chan *Message, channelBufSize)
	doneCh := make(chan bool)

	return &Client{maxId, ws, server, ch, doneCh, &UserInfo{}, "Unknown", "", false, db}
}

func (c *Client) StringLogin() string {
	return string(c.userInfo.Login)
}

func (c *Client) Conn() *websocket.Conn {
	return c.ws
}

func (c *Client) Write(msg *Message) {
	select {
	case c.ch <- msg:
	default:
		c.server.Del(c)
		err := fmt.Errorf("client %d is disconnected.", c.id)
		c.server.Err(err)
	}
}

func (c *Client) Done() {
	c.doneCh <- true
}

// Listen Write and Read request via chanel
func (c *Client) Listen() {
	go c.listenWrite()
	c.listenRead()
}

// Listen write request via chanel
func (c *Client) listenWrite() {
	log.Println("Listening write to client")
	for {
		select {

		// send message to the client
		case msg := <-c.ch:
			//log.Println("Send:", msg.Type)
			//log.Println("Send:", msg.Text)
			websocket.WriteJSON(c.ws, msg.Text)

		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

// Listen read request via chanel
func (c *Client) listenRead() {
	log.Println("Listening read from client")
	for {
		select {

		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
			var msg map[string]interface{}
			//err := websocket.JSON.Receive(c.ws, &msg)
			err := websocket.ReadJSON(c.ws, &msg)
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				c.ws.Close()
				c.doneCh <- true
				c.server.Err(err)
			} else {
				var ars []*APIReturn
				log.Println(msg)
				ParseAPI(c, &msg, &ars)
				//c.server.SendAll(&msg)
				for _, ar := range(ars) {
					switch ar.Type {
						default:
							c.server.SendQuery(ParseQuery(c, ar))
					}
				}
				ars = nil
			}
		}
	}
}