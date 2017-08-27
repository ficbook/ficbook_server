package chat

type Room struct {
	ID int `json:"id"`
	Name string `json:"name" gorm:"not null"`
	Topic string `json:"topic"`
	About string `json:"about" gorm:"type:text(500); not null" sql:"DEFAULT:'Unknown'"`
	Type string `json:"type" gorm:"not null"`
	UUID string `json:"uuid" gorm:"not null"`
	Power int `json:"-" gorm:"not null" sql:"DEFAULT:0"`
	Users map[int]*Client `json:"users" sql:"-"`
	LenUsers int `json:"count_users" sql:"-"`
}

//NewRoom returns the address of the room
func NewRoom(id int, name string, topic string, about, type_room string, UUID string, power int) *Room {
	return &Room{
		id,
		name,
		topic,
		about,
		type_room,
		UUID,
		power,
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
		0,
		make(map[int]*Client),
		0,
	}
}

func (r *Room) TableName() string {
	return "chat_rooms"
}
