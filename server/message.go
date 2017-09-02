package server

type Message struct {
	Client *Client
	Message *map[string]interface{}
}

/*
func NewChatMessageJSON(login string, message string, timestamp time.Time) ChatMessageJSON {
	return ChatMessageJSON{
		login,
		message,
		timestamp.UnixNano() / 1000000,
	}
}*/