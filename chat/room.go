package chat

type Room struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Topic string `json:"topic"`
	About string `json:"about"`
	UUID string `json:"uuid"`
	Users []*Client `json:"users"`
	LenUsers int `json:"count_users"`
}

//NewRoom returns the address of the room
func NewRoom(id int, name string, topic string, about, UUID string) *Room {
	return &Room{
		id,
		name,
		topic,
		about,
		UUID,
		[]*Client{},
		0,
	}
}

//CreateRoom returns the address of the created room
func CreateRoom(id int, name string, topic string, about string) *Room {
	return &Room{
		id,
		name,
		topic,
		about,
		NewUUID(),
		[]*Client{},
		0,
	}
}