package chat

//APIReturn - the structure returned by the ParseAPI function
type APIReturn struct {
	Type string
	Text string
}

type InfoQuery struct {
	Client *Client
	ApiReturn *APIReturn
}