package chat

import (
	"time"
	"strings"
	"strconv"
)


//ParseAPI works with the JSON requests
func ParseAPI(client *Client, msg *map[string]interface{}, apiReturn *APIReturn) {
	type_msg, ok := (*msg)["type"]
	if ok {
		*apiReturn = APIReturn{"MESSAGE", (*msg)["type"].(string), nil, nil}
		switch type_msg.(string) {
			case "autorize":
				isAuth := Authorization((*msg)["login"].(string), (*msg)["password"].(string))
				if isAuth {
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
					*apiReturn = APIReturn{"AUTH_OK", GetJSONUserInfo(userInfo.Login, userInfo.Password, userInfo.Power), nil, nil}
				} else {
					*apiReturn = APIReturn{"AUTH_ERROR", `{"type": "autorization", "result": "falled" ,"error": "erro"}`, nil, nil}
				}
			case "rooms":
				action_msg, _ := (*msg)["action"]
				switch action_msg {
					case "get":
						mapInterface := make(map[string]interface{})
						mapInterface["type"] = "rooms"
						//mapInterface["list"] = client.server.rooms
						var rooms []Room
						for _, room := range(client.server.rooms) {
							rooms = append(rooms, *room)
						}
						mapInterface["list"] = rooms
						*apiReturn = APIReturn{"ROOMS_GET", "", &mapInterface, nil}
				}
			case "room":
				action_msg, _ := (*msg)["action"]
				switch action_msg {
					case "join":
						room := client.server.GetSpecialRoomByName((*msg)["room_name"].(string))
						*apiReturn = APIReturn{"ROOM_JOIN", `{"type":"room","object":"about","room_name":"`+room.Name+`","about":"`+room.About+`"}`, nil, &ReturnVariable{nil, 0, room.Name}}
				}
			case "chat":
				action_msg, _ := (*msg)["action"]
				subject, _ := (*msg)["subject"]
				switch action_msg {
					case "get":
						switch subject {
							case "history":
								timestamp, _ := (*msg)["timestamp"]
								//tN := (time.Now().UnixNano() / 1000000) + 10800
								//client.server.db.Exec("SELECT current_timestamp();", tN)
								var messageSQL []*ChatMessageSQL
								//tN := time.Now().Local().Add(time.Hour * time.Duration(3) + time.Minute * time.Duration(0) + time.Second * time.Duration(0))
								//client.server.db.Table("chat_message_all").Order("timestamp desc").Where("timestamp BETWEEN ? AND ?", timestamp, tN).Find(&messageSQL).Limit(20)
								client.server.db.Table("chat_message_all").Order("timestamp desc").Where("timestamp <= ?", timestamp).Where("room_uuid = ?", client.server.GetSpecialRoomByName((*msg)["room_name"].(string)).UUID).Find(&messageSQL).Limit(20)
								var messageJSON []ChatMessageJSON
								for _, mes := range(messageSQL) {
									messageJSON = append(messageJSON, NewChatMessageJSON(mes.Login, mes.Message, mes.Timestamp))
								}
								//vv := ParseMessageQuery(client, messageJSON, apiReturn)
								*apiReturn = APIReturn{"CHAT_GET_HISTORY", "", nil, &ReturnVariable{&messageJSON, 0, ""}}
							case "participants":
								room := client.server.GetSpecialRoomByName((*msg)["room_name"].(string))
								messageJSON := make(map[string]interface{})
								messageJSON["type"] = "room"
								messageJSON["object"] = "about"
								messageJSON["room_name"] = room.Name
								messageJSON["about"] = room.About
								*apiReturn = APIReturn{"ROOM_GET_PARTICIPANTS", "", &messageJSON, nil}
						}
					case "send":
						switch subject {
							case "message":
								isCommand := false
								endMessage := (*msg)["message"].(string)
								if endMessage[0] == '!' {
									isCommand = true
								}
								room := client.server.GetSpecialRoomByName((*msg)["room_name"].(string))								
								if strings.Contains(room.Type, "system") && client.userInfo.Power < 100 {
									*apiReturn = *CreateCustomEvent("CHAT_SEND_MESSAGE", "You do not have permission to post in this room")
								} else {
									returnVariable := ReturnVariable{nil, 7777, ""}
									mapInterface := CreateInterfaceMessage((*msg)["room_name"].(string), "", "")
									if strings.Contains(room.Type, "system") || isCommand {
										(*mapInterface)["user"] = "Ficbook Chat Message"
									} else {
										(*mapInterface)["user"] = client.login
									}
									if isCommand {
										messages := strings.Split(endMessage, " ")
										returnVariable = ReturnVariable{nil, 0, ""}
										switch messages[0][1:] {
											default:
												endMessage = "This command does not exist. Enter !help to view commands"
											case "help":
												endMessage = "!test - Testing the command\n!refresh - Refresh"
											case "test":
												endMessage = "Test message!"
											case "refresh":
												if client.userInfo.Power < 1000 {
													endMessage = "You do not have permission to use this command"
												} else {
													endMessage = "refresh:\n\trooms"
													if len(messages) > 1 {
														switch messages[1] {
															case "rooms":
																client.server.RefreshRoom()
																endMessage = ""
															}
														}														
												}
											case "test2":
												room := client.server.GetSpecialRoomByName((*msg)["room_name"].(string))
												GetLoginUsers(room)
										}
									}
									(*mapInterface)["message"] = endMessage
									if !isCommand {
										client.server.db.Table("chat_message_all").Create(&ChatMessageSQL{
											Login: (*mapInterface)["user"].(string),
											Message: endMessage,
											Timestamp: time.Now(),
											RoomUUID: room.UUID,
										})
									}
									*apiReturn = APIReturn{"CHAT_SEND_MESSAGE", "", mapInterface, &returnVariable}
								}
						}
				}
		}
	} else {
		*apiReturn = APIReturn{"ERROR", `{"type": "Error", "result": "falled", "error": "Missing type key"}`, nil, nil}
	}
}

//ParseQuery returns a InfoQuery with data to work on the server
func ParseQuery(client *Client, apiReturn *APIReturn) *InfoQuery {
	return &InfoQuery{
		client,
		apiReturn,
	}
}

func ParseMessageQuery(client *Client, messages *[]ChatMessageJSON, apiReturn *APIReturn) *InfoQuery {
	mapInterface := make(map[string]interface{})
	mapInterface["type"] = "history"
	mapInterface["name"] = client.roomUUID
	mapInterface["messages"] = *messages
	if len(*messages) == 0 {
		mapInterface["messages"] = []int{}
	}
	return &InfoQuery{
		client,
		&APIReturn{
			"MESSAGE_SEND",
			"",
			&mapInterface,
			nil,
		},
	}
}

func CreateCustomEvent(typeAPI string, message string) *APIReturn {
	mapInterface := make(map[string]interface{})
	mapInterface["type"] = "event"
	mapInterface["action"] = "custom"
	mapInterface["message"] = message
	return &APIReturn{typeAPI, "", &mapInterface, nil}
}

func CreateEventUsersInfo(client *Client, action string, room *Room, code int) *APIReturn {
	if strings.Contains(action, "join") {
		room.Users = append(room.Users, client)
		room.LenUsers++
	} else if strings.Contains(action, "leave") {
		room.LenUsers--		
		RemoveAt(client, room)
	}
	mapInterface := make(map[string]interface{})
	mapInterface["type"] = "event"
	mapInterface["action"] = action
	mapInterface["users_count"] = room.LenUsers
	mapInterface["user_name"] = string(client.userInfo.Login)
	mapInterface["room_name"] = room.Name
	return &APIReturn{"CUSTOM_EVENT_COUNT_USERS", "", &mapInterface, &ReturnVariable{nil, code, room.UUID}}
}

func (client *Client) CreateEventUsersList(room *Room, code int) *APIReturn {
	mapInterface := make(map[string]interface{})
	mapInterface["type"] = "chat"
	mapInterface["action"] = "get"
	mapInterface["object"] = "participants"
	mapInterface["room_name"] = room.Name
	mapInterface["participants"] = *GetLoginUsers(room)
	return &APIReturn{"CHAT_GET_PARTICIPANTS_LOGINS", "", &mapInterface, &ReturnVariable{nil, code, room.UUID}}
}

func CreateInterfaceMessage(room_name string, user string, message string) *map[string]interface{} {
	if strings.Contains(user, "") {
		user = "Ficbook Chat Bot"
	}
	mapInterface := make(map[string]interface{})
	mapInterface["type"] = "chat"
	mapInterface["object"] = "message"
	mapInterface["time"] = time.Now().UnixNano() / 1000000
	mapInterface["room_name"] = room_name
	mapInterface["user"] = user
	mapInterface["message"] = message
	return &mapInterface
}

func GetJSONUserInfo(login []byte, password []byte, power int) string {
	return `{"type":"status", "action": "authorization", "status": "success", "power":` + strconv.Itoa(power) + `, "result": "ok", "login": "` + string(login) + `","password": "` + string(password) + `"}`
}

func GetLoginUsers(room *Room) *[]string {
	var logins []string
	for _, r := range(room.Users) {
		logins = append(logins, string(r.userInfo.Login))
	}
	if len(logins) == 0 {
		logins = []string{}
	}
	return &logins
}