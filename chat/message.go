package chat

type Message struct {
	Type string `json:"type"`
	Text map[string]interface{} `json:"data"`
}

func (self *Message) String() string {
	return "Send Message [" + self.Type + "]"
}
