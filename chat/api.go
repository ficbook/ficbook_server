package chat

import "strings"

//ParseAPI works with the JSON requests
func ParseAPI(msg *map[string]interface{}, apiReturn *APIReturn) {
	type_msg, ok := (*msg)["type"]
	if ok == false {
		*apiReturn = APIReturn{"ERROR", `Missing "type" key`}
	} else {
		*apiReturn = APIReturn{"MESSAGE", (*msg)["type"].(string)}
	}
	if strings.Contains(type_msg.(string), "autorize") || strings.Contains(type_msg.(string), "authorize") {
		isAuth := Authorization((*msg)["login"].(string), (*msg)["password"].(string))
		if isAuth {
			*apiReturn = APIReturn{"AUTH_OK", `{"type": "autorization", "result": "ok", "login": "` + (*msg)["login"].(string) + `","token": "token"}`}
		} else {
			*apiReturn = APIReturn{"AUTH_ERROR", `{"type": "autorization", "result": "falled" ,"error": "erro"}`}
		}
	}
}

//ParseQuery returns a InfoQuery with data to work on the server
func ParseQuery(client *Client, apiReturn *APIReturn) InfoQuery {
	return InfoQuery{
		client,
		apiReturn,
	}
}