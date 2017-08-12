package chat

import (
	"fmt"
	"io"
	"log"
	"golang.org/x/net/websocket"
	"github.com/jinzhu/gorm"
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
			if c.isAuth {
				websocket.JSON.Send(c.ws, msg.Text)
			} else {
				intf := make(map[string]interface{})
				intf["type"] = "error"
				intf["error"] = "You are not authorized"
				websocket.JSON.Send(c.ws, intf)
			}

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
			err := websocket.JSON.Receive(c.ws, &msg)
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				c.server.Err(err)
			} else {
				var ar APIReturn
				log.Println(msg)
				ParseAPI(c, &msg, &ar)
				switch ar.Type {
					default:
						vv := ParseQuery(c, &ar)
						c.server.SendQuery(vv)
					case "AUTH_OK":
						(*c).isAuth = true
						(*c).login = msg["login"].(string)
						vv := ParseQuery(c, &ar)
						c.server.SendQuery(vv)
					case "ROOM_JOIN":
						if len(c.roomUUID) > 1 {
							localRoom := c.server.GetSpecialRoomByName(c.roomUUID)
							vv := ParseQuery(c, CreateEventUsersInfo(c, "leave", localRoom, 3535))
							c.server.SendQuery(vv)
						}
						room := c.server.GetSpecialRoomByName(msg["room_name"].(string))
						log.Print(room)
						vv := ParseQuery(c, CreateEventUsersInfo(c, "join", room, 3535))
						c.server.SendQuery(vv)
						c.roomUUID = room.Name
						vv = ParseQuery(c, c.CreateEventUsersList(room, 3535))
						c.server.SendQuery(vv)
						mesJSON := make(map[string]interface{})
						mesJSON["type"] = "event"
						mesJSON["action"] = "room"
						mesJSON["object"] = "about"
						mesJSON["about"] = room.About
						vv = ParseQuery(c, &APIReturn{"USER_ROOM_ABOUT", "", &mesJSON, nil})
						c.server.SendQuery(vv)
						vv = ParseQuery(c, &ar)
						c.server.SendQuery(vv)
						c.roomUUID = vv.ApiReturn.ReturnVariable.string
						var messageSQL []*ChatMessageSQL
						c.server.db.Table("chat_message_all").Where("room_uuid = ?", room.UUID).Order("id desc").Find(&messageSQL).Limit(10)
						var messageJSON []ChatMessageJSON
						for _, mes := range(messageSQL) {
							messageJSON = append(messageJSON, NewChatMessageJSON(mes.Login, mes.Message, mes.Timestamp))
						}
						vv = ParseMessageQuery(c, &messageJSON, &ar)
						c.server.SendQuery(vv)
					case "CHAT_GET_HISTORY":
						vv := ParseMessageQuery(c, ar.ReturnVariable.ChatMessageJson, &ar)
						c.server.SendQuery(vv)
					case "ROOM_GET_PARTICIPANTS":
						vv := ParseQuery(c, &ar)
						c.server.SendQuery(vv)
						room := c.server.GetSpecialRoomByName(msg["room_name"].(string))
						vv = ParseQuery(c, c.CreateEventUsersList(room, 3535))
						c.server.SendQuery(vv)
				}
			}
		}
	}
}
