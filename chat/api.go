package chat

import (
	"strconv"
	"time"
	"strings"
	"github.com/gorilla/websocket"
)

//ParseAPI works with the JSON requests
func ParseAPI(client *Client, msg *map[string]interface{}, mapAPIReturn *[]*APIReturn) {
	typeMessage, ok := (*msg)["type"]
	if ok {
		switch typeMessage.(string) {
		case "autorize":
			isAuth := Authorization((*msg)["login"].(string), (*msg)["password"].(string))
			if isAuth {
				isBan := false
				var ban Ban
				client.server.db.Where("login_banned = ?", (*msg)["login"].(string)).First(&ban)
				if !ban.TimeBan.IsZero() {
					if time.Now().Nanosecond() >= ban.TimeExpired.Nanosecond() {
						localBan := Ban{LoginBanned:[]byte((*msg)["login"].(string))}
						client.server.db.Delete(&localBan)
					} else {
						isBan = true
					}
				}
				if !isBan {
					var userInfo UserInfo
					client.server.db.Where("login = ?", (*msg)["login"].(string)).First(&userInfo)
					if len(userInfo.Login) == 0 {
						userInfo = UserInfo{
							[]byte((*msg)["login"].(string)),
							[]byte((*msg)["password"].(string)),
							0,
							time.Now(),
							time.Now(),
						}
						client.server.db.Create(&userInfo)
					}
					(*client).userInfo = &userInfo
					(*client).isAuth = true
					returnMap := NewMap()
					GetMapUserInfo(returnMap, userInfo.Login, userInfo.Password, userInfo.Power)
					*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("AUTH_OK", returnMap, nil))
				} else {
					textMessage := "You are banned by " + string(ban.LoginBanning) + "\nReason: " + ban.Reason
					client.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, textMessage))
					client.Done()
				}
			} else {
				returnMap := NewMap()
				GetMapAuthError(returnMap)
				*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("AUTH_ERROR", returnMap, nil))
			}
		case "rooms":
			actionMessage, _ := (*msg)["action"]
			switch actionMessage {
			case "get":
				returnMap := NewMap()
				GetMapListRooms(returnMap, client.server.rooms)
				*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("ROOMS_GET", returnMap, nil))
			}
		case "room":
			actionMessage, _ := (*msg)["action"]
			switch actionMessage {
			case "join":
				if len(client.roomUUID) > 0 {
					localRoom := client.server.GetSpecialRoomByName(client.roomUUID)
					if localRoom.LenUsers > 0 {
						returnMap := NewMap()
						localRoom.LenUsers--
						localRoom.RemoveAt(client.id)
						GetMapEventUserCount(returnMap, string(client.userInfo.Login), "leave", localRoom.Name, localRoom.LenUsers)
						*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("ROOM_LEAVE", returnMap, NewReturnVariableRoom(localRoom, 35)))
					}
				}
				returnMap := NewMap()
				room := client.server.GetSpecialRoomByName((*msg)["room_name"].(string))
				GetMapRoomJoin(returnMap, room.Name, room.About)
				*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("ROOM_ABOUT", returnMap, NewReturnVariableString(room.Name, 35)))

				client.roomUUID = room.Name
				room.LenUsers++
				room.Users[client.id] = client
				returnMap = NewMap()
				GetMapEventUserCount(returnMap, string(client.userInfo.Login), "join", room.Name, room.LenUsers)
				*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("ROOM_JOIN", returnMap, NewReturnVariableRoom(room, 35)))

				returnMap = NewMap()
				GetMapRoomAbout(returnMap, room.Name, room.About)
				*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("ROOM_ABOUT", returnMap, nil))

				var messageSQL []ChatMessageSQL
				var messageJSON []ChatMessageJSON
				returnMap = NewMap()
				client.server.db.Table("chat_message_all").Where("room_uuid = ?", room.UUID).Order("id desc").Find(&messageSQL).Limit(10)
				for _, mes := range messageSQL {
					messageJSON = append(messageJSON, NewChatMessageJSON(mes.Login, mes.Message, mes.Timestamp))
				}
				GetMapHistoryMessages(returnMap, room.Name, &messageJSON)
				*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("ROOM_HISTORY", returnMap, nil))

				returnMap = NewMap()
				GetMapEventUserList(returnMap, room.Name, GetLoginUsers(room.Users))
				*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("ROOM_USERS", returnMap, nil))
			}
		case "chat":
			actionMessage, _ := (*msg)["action"]
			subject, _ := (*msg)["subject"]
			switch actionMessage {
			case "get":
				switch subject {
				case "history":
					timestamp, _ := (*msg)["timestamp"]
					var messageSQL []*ChatMessageSQL
					var messageJSON []ChatMessageJSON
					returnMap := NewMap()
					client.server.db.Table("chat_message_all").Order("timestamp desc").Where("timestamp <= ?", timestamp).Where("room_uuid = ?", client.server.GetSpecialRoomByName((*msg)["room_name"].(string)).UUID).Find(&messageSQL).Limit(20)
					for _, mes := range messageSQL {
						messageJSON = append(messageJSON, NewChatMessageJSON(mes.Login, mes.Message, mes.Timestamp))
					}
					GetMapHistoryMessages(returnMap, (*msg)["room_name"].(string), &messageJSON)
					*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("CHAT_GET_HISTORY", returnMap, nil))
				case "participants":
					room := client.server.GetSpecialRoomByName((*msg)["room_name"].(string))
					returnMap := NewMap()
					GetMapRoomAbout(returnMap, room.Name, room.About)
					*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("ROOM_GET_PARTICIPANTS", returnMap, nil))

					GetMapEventUserList(returnMap, room.Name, GetLoginUsers(room.Users))
					*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("ROOM_USERS", returnMap, nil))
				}
			case "send":
					switch subject {
						case "message":
							isCommand := false
							userName := string(client.userInfo.Login)
							endMessage := (*msg)["message"].(string)
							if endMessage[0] == '!' {
								isCommand = true
							}
							room := client.server.GetSpecialRoomByName((*msg)["room_name"].(string))								
							if room.Type == "system" && client.userInfo.Power < 10000 {
								returnMap := NewMap()
								GetMapCustomEvent(returnMap, "You do not have permission to post in this room")
								*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("CHAT_SEND_MESSAGE", returnMap, nil))
							} else {
								returnMap := NewMap()
								GetMapCreateMessage(returnMap, (*msg)["room_name"].(string), "", "")
								if room.Type == "system" || isCommand {
									userName = "System message"
								}								
								if isCommand {
									messages := strings.Split(endMessage, " ")
									switch messages[0][1:] {
										default:
											endMessage = "This command does not exist. Enter !help to view commands"
										case "help":
											endMessage = "!test - Testing the command\n!refresh - Refresh"
										case "test":
											endMessage = "Test message!"
										case "refresh":
											if client.userInfo.Power < 10000 {
												endMessage = "You do not have permission to use this command"
											} else {
												endMessage = "refresh:\n\trooms"
												if len(messages) > 1 {
													switch messages[1] {
														case "rooms":
															client.server.RefreshRoom()
															endMessage = "Комнаты обновлены"
														}
													}														
											}
									}
								}
								(*returnMap)["user"] = userName
								(*returnMap)["message"] = endMessage
								if !isCommand {
									client.server.db.Table("chat_message_all").Create(&ChatMessageSQL{
										Login: userName,
										Message: endMessage,
										Timestamp: time.Now(),
										RoomUUID: room.UUID,
									})
									*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("CHAT_SEND_MESSAGE", returnMap, NewReturnVariableRoom(room, 35)))
								} else {
									*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("CHAT_SEND_ADM_MESSAGE", returnMap, nil))
								}
							}
						}
			}
		case "administration":
			actionMessage, _ := (*msg)["action"]
			switch actionMessage {
				case "get":
					objectMessage, _ := (*msg)["object"]
					switch objectMessage {
						case "bans":
							var bans []Ban
							returnMap := NewMap()
							client.server.db.Find(&bans)
							GetMapListBans(returnMap, &bans)
							*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("ADM_GET_BANS", returnMap, nil))
					}
				case "kik":
					localClient, isSearch := client.server.SearchUser((*msg)["user_name"].(string))
					if isSearch {
						if client.userInfo.Power > localClient.userInfo.Power {
							room := client.server.GetSpecialRoomByName(localClient.roomUUID)
							if room != nil {
								room.LenUsers--
								returnMap := NewMap()
								GetMapUserClosed(returnMap, room.LenUsers, localClient.StringLogin(), room.Name)
								*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("ADM_CLOSE_INFO", returnMap, NewReturnVariableRoom(room, 35)))
							}

							returnMap := NewMap()
							GetMapCustomEvent(returnMap, "You kicked " + (*msg)["user_name"].(string),)
							*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("CHAT_CUSTOM_MESSAGE", returnMap, nil))

							textMessage := "You are kicked by " + client.StringLogin() + "\nReason: " + (*msg)["message"].(string)
							localClient.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, textMessage))
							localClient.Done()
						}
					}
				case "ban":
					localClient, isSearch := client.server.SearchUser((*msg)["user_name"].(string))
					if isSearch {
						if client.userInfo.Power > localClient.userInfo.Power {
							ban := Ban{
								localClient.userInfo.Login,
								client.userInfo.Login,
								(*msg)["message"].(string),
								time.Now(),
								time.Now().Add(time.Duration(60) * time.Minute),
								//time.Now().Add(time.Duration((*msg)["duration"].(float64)) * time.Nanosecond),
							}
							client.server.db.Create(&ban)

							room := client.server.GetSpecialRoomByName(localClient.roomUUID)
							if room != nil {
								room.LenUsers--
								returnMap := NewMap()
								GetMapUserClosed(returnMap, room.LenUsers, localClient.StringLogin(), room.Name)
								*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("ADM_CLOSE_INFO", returnMap, NewReturnVariableRoom(room, 35)))
							}

							returnMap := NewMap()
							GetMapCustomEvent(returnMap, "You banned " + (*msg)["user_name"].(string),)
							*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("CHAT_CUSTOM_MESSAGE", returnMap, nil))

							textMessage := "You are banned by " + client.StringLogin() + "\nReason: " + (*msg)["message"].(string)
							localClient.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, textMessage))
							localClient.Done()
						}
					}
				case "create":
					objectMessage, _ := (*msg)["object"]
					switch objectMessage {
						case "room":
							room := CreateRoom(0, (*msg)["name"].(string), "", "Unknown", "public")
							client.server.db.Create(room)
							client.server.db.Where("uuid = ?", (*room).UUID).First(room)
							client.server.rooms[(*room).ID] = room
							
							returnMap := NewMap()
							GetMapCustomEvent(returnMap, "You have created a room with the name " + (*msg)["name"].(string),)
							*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("ADM_INFO_CREATE_ROOM", returnMap, nil))
					}
			}
		}
	} else {
		returnMap := NewMap()
		GetMapTypeError(returnMap)
		*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("ERROR", returnMap, nil))
	}
}

//ParseQuery returns a InfoQuery with data to work on the server
func ParseQuery(client *Client, apiReturn *APIReturn) *InfoQuery {
	return &InfoQuery{
		client,
		apiReturn,
	}
}

func NewMap() *map[string]interface{} {
	returnMap := make(map[string]interface{})
	return &returnMap
}

func GetMapTypeError(returnMap *map[string]interface{}) {
	(*returnMap)["type"] = "error"
	(*returnMap)["result"] = "falled"
	(*returnMap)["error"] = "Missing type key"
}

func GetMapUserInfo(returnMap *map[string]interface{}, login []byte, password []byte, power int) {
	(*returnMap)["type"] = "status"
	(*returnMap)["action"] = "authorization"
	(*returnMap)["status"] = "success"
	(*returnMap)["power"] = strconv.Itoa(power)
	(*returnMap)["result"] = "ok"
	(*returnMap)["login"] = string(login)
	(*returnMap)["password"] = string(password)
}

func GetMapAuthError(returnMap *map[string]interface{}) {
	(*returnMap)["type"] = "autorization"
	(*returnMap)["result"] = "falled"
	(*returnMap)["error"] = "erro"
}

func GetMapListRooms(returnMap *map[string]interface{}, rooms map[int]*Room) {
	var returnRooms []Room
	for _, room := range rooms {
		returnRooms = append(returnRooms, *room)
	}
	(*returnMap)["type"] = "rooms"
	(*returnMap)["list"] = returnRooms
}

func GetMapRoomJoin(returnMap *map[string]interface{}, roomName string, roomAbout string) {
	(*returnMap)["type"] = "room"
	(*returnMap)["object"] = "about"
	(*returnMap)["room_name"] = roomName
	(*returnMap)["about"] = roomAbout
}

func GetMapEventUserCount(returnMap *map[string]interface{}, userLogin string, action string, roomName string, roomLenUsers int) {
	(*returnMap)["type"] = "event"
	(*returnMap)["action"] = action
	(*returnMap)["users_count"] = roomLenUsers
	(*returnMap)["user_name"] = userLogin
	(*returnMap)["room_name"] = roomName
}

func GetMapRoomAbout(returnMap *map[string]interface{}, roomName string, roomAbout string) {
	(*returnMap)["type"] = "room"
	(*returnMap)["object"] = "about"
	(*returnMap)["room_name"] = roomName
	(*returnMap)["about"] = roomAbout
}

func GetMapRoomAbout2(returnMap *map[string]interface{}, roomAbout string) {
	(*returnMap)["type"] = "event"
	(*returnMap)["action"] = "room"
	(*returnMap)["object"] = "about"
	(*returnMap)["about"] = roomAbout
}

func GetMapHistoryMessages(returnMap *map[string]interface{}, roomName string, messageJSON *[]ChatMessageJSON) {
	(*returnMap)["type"] = "history"
	(*returnMap)["name"] = roomName
	if len(*messageJSON) == 0 {
		(*returnMap)["messages"] = []string{}
	} else {
		(*returnMap)["messages"] = *messageJSON
	}
}

func GetMapEventUserList(returnMap *map[string]interface{}, roomName string, users *[]string) {
	(*returnMap)["type"] = "chat"
	(*returnMap)["action"] = "get"
	(*returnMap)["object"] = "participants"
	(*returnMap)["room_name"] = roomName
	(*returnMap)["participants"] = *users
}

func GetMapCustomEvent(returnMap *map[string]interface{}, messageText string) {
	(*returnMap)["type"] = "event"
	(*returnMap)["action"] = "custom"
	(*returnMap)["message"] = messageText
}

func GetMapCreateMessage(returnMap *map[string]interface{}, roomName string, userName string, messageText string) {
	if strings.Contains(userName, "") {
		userName = "System message"
	}
	(*returnMap)["type"] = "chat"
	(*returnMap)["object"] = "message"
	(*returnMap)["time"] = time.Now().UnixNano() / 1000000
	(*returnMap)["room_name"] = roomName
	(*returnMap)["user"] = userName
	(*returnMap)["message"] = messageText
}

func GetMapListBans(returnMap *map[string]interface{}, bans *[]Ban) {
	(*returnMap)["type"] = "list"
	(*returnMap)["object"] = "bans"
	if len(*bans) == 0 {
		(*returnMap)["list"] = []string{}
	} else {
		(*returnMap)["list"] = *bans
	}
}

func GetMapUserClosed(returnMap *map[string]interface{}, userCount int, userName string, roomName string) { 
	(*returnMap)["type"] = "event"
	(*returnMap)["action"] = "kiked"
	(*returnMap)["users_count"] = userCount
	(*returnMap)["user_name"] = userName
	(*returnMap)["room_name"] = roomName
}