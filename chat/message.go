package chat

type Message struct {
	Login string `json:"login"`
	Text string `json:"text"`
}

func (self *Message) String() string {
	return self.Login + ": " + self.Text
}
