package chat

import (
	"sort"
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
							Login: []byte((*msg)["login"].(string)),
							Password: []byte((*msg)["password"].(string)),
							Power: 0,
							DateReg: time.Now(),
							DateVisit: time.Now(),
						}
						client.server.db.Create(&userInfo)
					}
					(*client).userInfo = &userInfo
					(*client).isAuth = true
					returnMap := NewMap()
					GetMapUserInfo(returnMap, userInfo.Login, userInfo.Password, userInfo.Power)
					*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("AUTH_OK", returnMap, nil))
				} else {
					textMessage := (*client.server.lang)["banned_info_1"] + string(ban.LoginBanning) + "\n" + (*client.server.lang)["banned_info_2"] + ban.Reason
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
				GetMapListRooms(returnMap, client.server.roomsList)
				*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("ROOMS_GET", returnMap, nil))
			}
		case "room":
			actionMessage, _ := (*msg)["action"]
			switch actionMessage {
				case "join":
					if client.room != nil {
						if client.room.LenUsers > 0 {
							returnMap := NewMap()
							client.room.LenUsers--
							client.room.RemoveAt(client.id)
							GetMapEventUserCount(returnMap, client.StringLogin(), "leave", client.room.Name, client.room.LenUsers)
							*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("ROOM_LEAVE", returnMap, NewReturnVariableRoom(client.room, 35)))
						}
					}
					returnMap := NewMap()
					room := client.server.GetSpecialRoomByName((*msg)["room_name"].(string))
					GetMapRoomJoin(returnMap, room.Name, room.About)
					*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("ROOM_ABOUT", returnMap, NewReturnVariableString(room.Name, 35)))

					client.room = room
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
				
				case "set":
					subject, _ := (*msg)["subject"]
					switch subject {
						case "about":
							room := client.server.GetSpecialRoomByName((*msg)["room_name"].(string))
							if room != nil {
								room.About = (*msg)["about"].(string)
								
								client.server.db.Save(room)

								returnMap := NewMap()
								GetMapCustomEvent(returnMap, (*client.server.lang)["change_about_room"])
								*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("ROOM_SET_ABOUT", returnMap, nil))

								returnMap = NewMap()
								GetMapRoomAbout(returnMap, room.Name, room.About)
								*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("ROOM_SET_ABOUT", returnMap, NewReturnVariableRoom(room, 35)))
							}
					}

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
							room := client.server.GetSpecialRoomByName((*msg)["room_name"].(string))
							if room != nil {
								if client.antiflood > 5 {
									returnMap := NewMap()
									GetMapCustomEvent(returnMap, client.StringLogin() + (*client.server.lang)["antiflood_info_1"])
									*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("CHAT_CUSTOM_MESSAGE", returnMap, NewReturnVariableRoom(room, 35)))
		
									client.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, (*client.server.lang)["antiflood_info_2"]))
									client.Done()
								}
								client.antiflood++
								isCommand := false
								userName := string(client.userInfo.Login)
								endMessage := (*msg)["message"].(string)
								if endMessage[0] == '!' {
									isCommand = true
								}						
								if room.Type == "system" && client.userInfo.Power < 10000 {
									returnMap := NewMap()
									GetMapCustomEvent(returnMap, (*client.server.lang)["dont_have_permission"])
									*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("CHAT_SEND_MESSAGE", returnMap, nil))
								} else {
									returnMap := NewMap()
									GetMapCreateMessage(returnMap, (*client.server.lang)["server_name"], (*msg)["room_name"].(string), "", "")
									if room.Type == "system" || isCommand {
										userName = (*client.server.lang)["server_name"]
									}								
									if isCommand {
										messages := strings.Split(endMessage, " ")
										switch messages[0][1:] {
											default:
												endMessage = (*client.server.lang)["command_not_exist"]
											case "help":
												endMessage = (*client.server.lang)["commands_1"] + "\n" + (*client.server.lang)["commands_2"]
												endMessage += "\n"+ (*client.server.lang)["commands_3"] + "\n" + (*client.server.lang)["commands_4"]
											case "test":
												endMessage = (*client.server.lang)["commands_result_1"]
											case "refresh":
												if client.userInfo.Power < 10000 {
													endMessage = (*client.server.lang)["dont_have_permission"]
												} else {
													endMessage = "refresh:\n\trooms\n\tusers"
													if len(messages) > 1 {
														switch messages[1] {
															case "rooms":
																client.server.RefreshRoom()
																endMessage = (*client.server.lang)["rooms_updated"]
															case "users":
																textMessage := client.StringLogin() + " kicked all server"
																for _, lRoom := range(client.server.rooms) {
																	for _, lUser := range(lRoom.Users) {
																		lUser.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, textMessage))
																		lUser.Done()
																	}
																}
																client.ClearBufferID()
															}

														}														
												}
											case "users":
												endMessage = (*client.server.lang)["commands_result_3"] + "\n\n"
												var intsList []int
												for k := range(client.room.Users) {
													intsList = append(intsList, k)
												}
												sort.Ints(intsList)
												for _, u := range(intsList) {
													endMessage += strconv.Itoa(client.room.Users[u].id) + " — " + client.room.Users[u].StringLogin() + "\n"
												}
											case "admins":
												endMessage = (*client.server.lang)["commands_result_4"] + "\n\n"
												var intsList []int
												for k := range(client.room.Users) {
													intsList = append(intsList, k)
												}
												sort.Ints(intsList)
												for _, u := range(intsList) {
													if client.room.Users[u].userInfo.Power <= 0 {
														continue
													}
													endMessage += strconv.Itoa(client.room.Users[u].id) + " — " + GetStringPrivilege(client.server.lang, client.room.Users[u].userInfo.Power) + " — " + client.room.Users[u].StringLogin() + "\n"
												}
											case "setadmin":
												if client.userInfo.Power >= 10000 {
													if len(messages) > 2 {
														idUser, _ := strconv.Atoi(messages[1])
														user, ok := client.room.Users[idUser]
														if ok {
															power := -1
															switch messages[2] {
																case "0":
																	power = 0
																case "1":
																	power = 100
																case "2":
																	power = 1000
																default:
																	endMessage = "!setadmin userID privilegeID\n0 - user\n1 - moder\n2 - admin"
															}
															if power > -1 {
																user.userInfo.Power = power
																client.server.db.Save(user.userInfo)

																endMessage = (*client.server.lang)["set_admin_result_1_1"] + user.StringLogin() + (*client.server.lang)["set_admin_result_1_2"] + GetStringPrivilege(client.server.lang, user.userInfo.Power)

																returnMap := NewMap()																
																GetMapCreateMessage(returnMap, "", user.room.Name, user.StringLogin(), client.StringLogin() + (*client.server.lang)["set_admin_result_2"] + GetStringPrivilege(client.server.lang, user.userInfo.Power))
																*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("CHAT_SEND_MESSAGE", returnMap, NewReturnVariableClient(user, 75)))
															}
														} else {
															endMessage = "User is not found!"
														}
													} else {
														endMessage = "!setadmin userID privilegeID\n0 - user\n1 - moder\n2 - admin"
													}
												} else {
													endMessage = (*client.server.lang)["dont_have_permission"]
												}
											case "rooms":
												endMessage = (*client.server.lang)["commands_result_5"] + "\n\n"
												for _, r := range(*client.server.roomsList) {
													endMessage += strconv.Itoa((*client.server.rooms[r.ID]).ID) + " — " + (*client.server.rooms[r.ID]).Name + "\n"
												}
											case "settype":
												if client.userInfo.Power >= 10000 {
													if len(messages) > 2 {
														idRoom, _ := strconv.Atoi(messages[1])
														if room, ok := client.server.rooms[idRoom]; ok {
															var (
																typeRoom string
																typeOk bool = true
															)
															switch messages[2] {
																case "0":
																	typeRoom = "public"
																case "1":
																	typeRoom = "system"
																default:
																	endMessage = "!settype roomID typeRoom\n0 - public\n1 - system"
																	typeOk = false
															}
															if typeOk {
																room.Type = typeRoom
																client.server.db.Save(room)
																endMessage = (*client.server.lang)["commands_result_6"]
															} else {
																endMessage = "!settype roomID typeRoom\n0 - public\n1 - system"
															}
														} else {
															endMessage = "Room is not found!"
														}
													} else {
														endMessage = "!settype roomID typeRoom\n0 - public\n1 - system"
													}
												} else {
													endMessage = (*client.server.lang)["dont_have_permission"]
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
							} else {
								returnMap := NewMap()
								GetMapCustomEvent(returnMap, (*client.server.lang)["room_is_deleted"])
								*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("CHAT_SEND_ERROR", returnMap, nil))
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
							if client.userInfo.Power >= 100 {
								var bans []Ban
								returnMap := NewMap()
								client.server.db.Find(&bans)
								GetMapListBans(returnMap, &bans)
								*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("ADM_GET_BANS", returnMap, nil))
							}
					}
				case "kik":
					localClient, isSearch := client.server.SearchUser((*msg)["user_name"].(string))
					if client.userInfo.Power >= 100 && isSearch {
						if client.userInfo.Power > localClient.userInfo.Power {
							if localClient.room != nil {
								localClient.room.LenUsers--
								returnMap := NewMap()
								GetMapUserClosed(returnMap, localClient.room.LenUsers, localClient.StringLogin(), localClient.room.Name)
								*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("ADM_CLOSE_INFO", returnMap, NewReturnVariableRoom(localClient.room, 35)))
							}

							client.server.db.Create(&AdminHistory{
								LoginAdmin: client.userInfo.Login,
								LoginUser: localClient.userInfo.Login,
								Action: "kick",
								Date: time.Now(),
							})

							returnMap := NewMap()
							GetMapCustomEvent(returnMap,(*client.server.lang)["kicked_info_1"] + (*msg)["user_name"].(string),)
							*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("CHAT_CUSTOM_MESSAGE", returnMap, nil))

							textMessage := (*client.server.lang)["kicked_info_2"] + client.StringLogin() + "\n" + (*client.server.lang)["kicked_info_3"] + (*msg)["message"].(string)
							localClient.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, textMessage))
							localClient.Done()
						}
					} else {
						returnMap := NewMap()
						GetMapCustomEvent(returnMap, (*client.server.lang)["dont_have_permission"])
						*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("ADM_INFO_DONT_HAVE_PERMISSION", returnMap, nil))
					}
				case "ban":
					localClient, isSearch := client.server.SearchUser((*msg)["user_name"].(string))
					if client.userInfo.Power >= 1000 && isSearch {
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

							client.server.db.Create(&AdminHistory{
								LoginAdmin: client.userInfo.Login,
								LoginUser: localClient.userInfo.Login,
								Action: "ban",
								Date: time.Now(),
							})

							if localClient.room != nil {
								localClient.room.LenUsers--
								returnMap := NewMap()
								GetMapUserClosed(returnMap, localClient.room.LenUsers, localClient.StringLogin(), localClient.room.Name)
								*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("ADM_CLOSE_INFO", returnMap, NewReturnVariableRoom(localClient.room, 35)))
							}

							returnMap := NewMap()
							GetMapCustomEvent(returnMap, (*client.server.lang)["banned_info_3"] + (*msg)["user_name"].(string),)
							*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("CHAT_CUSTOM_MESSAGE", returnMap, nil))

							textMessage := (*client.server.lang)["banned_info_1"] + client.StringLogin() + "\n" + (*client.server.lang)["banned_info_2"] + (*msg)["message"].(string)
							localClient.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, textMessage))
							localClient.Done()
						}
					} else {
						returnMap := NewMap()
						GetMapCustomEvent(returnMap, (*client.server.lang)["dont_have_permission"])
						*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("ADM_INFO_DONT_HAVE_PERMISSION", returnMap, nil))
					}
				case "create":
					objectMessage, _ := (*msg)["object"]
					switch objectMessage {
						case "room":
							if client.userInfo.Power >= 100 {
								room := CreateRoom(0, (*msg)["name"].(string), "", "Unknown", "public")
								client.server.db.Create(room)
								client.server.db.Where("uuid = ?", (*room).UUID).First(room)
								client.server.rooms[(*room).ID] = room
								
								client.server.UpdateListRooms()

								client.server.db.Create(&AdminHistory{
									LoginAdmin: client.userInfo.Login,
									LoginUser: []byte("-"),
									Action: "create_room",
									Date: time.Now(),
								})

								returnMap := NewMap()
								GetMapCustomEvent(returnMap, (*client.server.lang)["created_room"] + (*msg)["name"].(string),)
								*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("ADM_INFO_CREATE_ROOM", returnMap, nil))
							} else {
								returnMap := NewMap()
								GetMapCustomEvent(returnMap, (*client.server.lang)["dont_have_permission"])
								*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("ADM_INFO_DONT_HAVE_PERMISSION", returnMap, nil))
							}
					}
				case "destroy":
					objectMessage, _ := (*msg)["object"]
					switch objectMessage {			
						case "room":
							room := client.server.GetSpecialRoomByName((*msg)["room_name"].(string))
							if room != nil {
								if client.userInfo.Power >= 1000 {
									client.server.db.Delete(room)

									delete(client.server.rooms, room.ID)

									client.server.UpdateListRooms()

									client.server.db.Create(&AdminHistory{
										LoginAdmin: client.userInfo.Login,
										LoginUser: []byte("-"),
										Action: "delete_room",
										Date: time.Now(),
									})

									returnMap := NewMap()
									GetMapCustomEvent(returnMap, (*client.server.lang)["room_deleted"])
									*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("ADM_INFO_DELETE_ROOM", returnMap, nil))
								} else {
									returnMap := NewMap()
									GetMapCustomEvent(returnMap, (*client.server.lang)["dont_have_permission"])
									*mapAPIReturn = append(*mapAPIReturn, NewAPIReturn("ADM_INFO_DONT_HAVE_PERMISSION", returnMap, nil))
								}
							}
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

func GetMapListRooms(returnMap *map[string]interface{}, rooms *[]*Room) {
	var returnRooms []Room
	for _, room := range(*rooms) {
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

func GetMapCreateMessage(returnMap *map[string]interface{}, systemName string, roomName string, userName string, messageText string) {
	if strings.Contains(userName, "") {
		userName = systemName
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

func GetStringPrivilege(lang *map[string]string, power int) string {
	if power >= 10000 {
		return (*lang)["superadmin"]
	}
	if power >= 1000 {
		return (*lang)["admin"]
	}
	if power >= 100 {
		return (*lang)["moder"]
	}
	return (*lang)["user"]
}