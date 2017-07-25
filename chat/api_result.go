package chat

//APIReturn - the structure returned by the ParseAPI function
type APIReturn struct {
	Type string `json:"type"`
	Text string `json:"text"`
	Interface *map[string]interface{}
	ReturnVariable *ReturnVariable
}

type InfoQuery struct {
	Client *Client
	ApiReturn *APIReturn
}

type ReturnVariable struct {
	ChatMessageJson *[]ChatMessageJSON
	code int
	string
}