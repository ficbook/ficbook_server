package chat

type Message struct {
	Login string `json:"type"`
	Text string `json:"data"`
}

func (self *Message) String() string {
	return self.Login + ": " + self.Text
}
