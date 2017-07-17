package chat

//APIReturn - the structure returned by the ParseAPI function
type APIReturn struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type InfoQuery struct {
	Client *Client
	ApiReturn *APIReturn
}