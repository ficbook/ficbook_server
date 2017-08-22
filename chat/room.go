package chat

type Room struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Topic string `json:"topic"`
	About string `json:"about"`
	Type string `json:"type"`
	UUID string `json:"uuid"`
	Users map[int]*Client `json:"users" sql:"-"`
	LenUsers int `json:"count_users" sql:"-"`
}

//NewRoom returns the address of the room
func NewRoom(id int, name string, topic string, about, type_room string, UUID string) *Room {
	return &Room{
		id,
		name,
		topic,
		about,
		type_room,
		UUID,
		make(map[int]*Client),
		0,
	}
}

//CreateRoom returns the address of the created room
func CreateRoom(id int, name string, topic string, about string, type_room string) *Room {
	return &Room{
		id,
		name,
		topic,
		about,
		type_room,
		NewUUID(),
		make(map[int]*Client),
		0,
	}
}