package chat

type Room struct {
	ID int `json:"id"`
	Name string `json:"name" gorm:"not null"`
	Topic string `json:"topic"`
	About string `json:"about" gorm:"type:text(500); not null"`
	Type string `json:"type" gorm:"not null"`
	UUID string `json:"uuid" gorm:"not null"`
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

func (r *Room) TableName() string {
	return "chat_rooms"
}
