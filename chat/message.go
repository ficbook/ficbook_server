package chat

type Message struct {
	Login string `json:"author"`
	Body   string `json:"body"`
}

func (self *Message) String() string {
	return self.Login + " says " + self.Body
}
