package server

type Room struct {
	ID int `gorm:"primary_key; not null; AUTO_INCREMENT" json:"id"`
	Name string `gorm:"not null" json:"name"`
	Members int `json:"members" sql:"-"`
}

func NewRoom(id int, name string) *Room {
	return &Room{
		id,
		name,
		0,
	}
}