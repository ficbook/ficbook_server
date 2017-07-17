package chat
//ParseAPI works with the JSON requests
func ParseAPI(msg *map[string]interface{}, apiReturn *APIReturn) {
	_, ok := (*msg)["type"]
	if ok == false {
		*apiReturn = APIReturn{"ERROR", `Missing "type" key`}
	} else {
		*apiReturn = APIReturn{"MESSAGE", (*msg)["type"].(string)}
	}
	_, ok = (*msg)["test"]
	if ok {
		*apiReturn = APIReturn{"TEST", (*msg)["test"].(string)}
	}
}

//ParseQuery returns a map [string] interface {} with data to work on the server
func ParseQuery(client *Client, apiReturn *APIReturn) InfoQuery {
	return InfoQuery{
		client,
		apiReturn,
	}
}