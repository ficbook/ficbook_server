package chat


//ParseAPI works with the JSON requests
func ParseAPI(client *Client, msg *map[string]interface{}, apiReturn *APIReturn) {
	type_msg, ok := (*msg)["type"]
	if ok {
		*apiReturn = APIReturn{"MESSAGE", (*msg)["type"].(string), nil}
		switch type_msg.(string) {
			case "autorize":
			case "authorize":
				isAuth := Authorization((*msg)["login"].(string), (*msg)["password"].(string))
				if isAuth {
					*apiReturn = APIReturn{"AUTH_OK", `{"type":"status", "action": "authorization", "status": "success", "power": 0, "result": "ok", "login": "` + (*msg)["login"].(string) + `","password": "` + GenerateToken() + `"}`, nil}
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
						*apiReturn = APIReturn{"GET_ROOMS", "", &mapInterface}
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