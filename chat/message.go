package chat

type Message struct {
	Type string `json:"type"`
	Text map[string]interface{} `json:"data"`
}

func (self *Message) String() string {
	//return self.Login + ": " + self.Text
	return "1"
}
