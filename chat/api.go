package chat

//ParseAPI works with the JSON requests
func ParseAPI(msg *map[string]interface{}, apiReturn *APIReturn) {
	_, ok := (*msg)["type"]
	if ok == false {
		*apiReturn = APIReturn{"ERROR", `Missing "type" key`}
	} else {
		*apiReturn = APIReturn{"MESSAGE", (*msg)["type"].(string)}
	}
}
