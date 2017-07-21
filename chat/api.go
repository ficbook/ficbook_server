package chat

import "fmt"
//ParseAPI works with the JSON requests
func ParseAPI(client *Client, msg *map[string]interface{}, apiReturn *APIReturn) {
	type_msg, ok := (*msg)["type"]
	if ok {
		*apiReturn = APIReturn{"MESSAGE", (*msg)["type"].(string), nil}
		switch type_msg.(string) {
			case "autorize":
				isAuth := Authorization((*msg)["login"].(string), (*msg)["password"].(string))
				if isAuth {
					*apiReturn = APIReturn{"AUTH_OK", `{"type":"status", "action": "authorization", "status": "success", "power": 0, "result": "ok", "login": "` + (*msg)["login"].(string) + `","password": "` + (*msg)["password"].(string) + `"}`, nil}
				} else {
					*apiReturn = APIReturn{"AUTH_ERROR", `{"type": "autorization", "result": "falled" ,"error": "erro"}`, nil}
				}
			case "rooms":
				action_msg, _ := (*msg)["action"]
				switch action_msg {
					case "get":
						mapInterface := make(map[string]interface{})
						mapInterface["type"] = "rooms"
						mapInterface["list"] = client.server.rooms
						*apiReturn = APIReturn{"ROOMS_GET", "", &mapInterface}
				}
			case "room":
				action_msg, _ := (*msg)["action"]
				switch action_msg {
					case "join":
						*apiReturn = APIReturn{"ROOM_JOIN", `{"type":"room","object":"about","room_name":"`+(*msg)["room_name"].(string)+`","about":"Unknown"}`, nil}
				}
			}
	} else {
		*apiReturn = APIReturn{"ERROR", `{"type": "Error", "result": "falled", "error": "Missing type key"}`, nil}
	}
}

//ParseQuery returns a InfoQuery with data to work on the server
func ParseQuery(client *Client, apiReturn *APIReturn) InfoQuery {
	return InfoQuery{
		client,
		apiReturn,
	}
}