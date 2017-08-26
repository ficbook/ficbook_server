package chat

import (
	"time"
)

//APIReturn - the structure returned by the ParseAPI function
type APIReturn struct {
	Type string `json:"type"`
	Interface *map[string]interface{}
	ReturnVariable *ReturnVariable
}

type InfoQuery struct {
	Client *Client
	ApiReturn *APIReturn
}

type ReturnVariable struct {
	ChatMessageJSON *[]ChatMessageJSON
	string
	int
	ReturnRoom *Room
}

type Ban struct {
	LoginBanned []byte `json:"login_banned" gorm:"not null"`
	LoginBanning []byte `json:"login_baning" gorm:"not null"`
	Reason string `json:"reason" gorm:"not null"`
	TimeBan time.Time `json:"-" gorm:"not null"`
	TimeExpired time.Time `json:"time_expired" gorm:"not null"`
}

func (b *Ban) TableName() string {
	return "bans_list"
}

func NewAPIReturn(typeAPI string, interfaceAPI *map[string]interface{}, returnVariable *ReturnVariable) *APIReturn {
	return &APIReturn{
		typeAPI,
		interfaceAPI,
		returnVariable,
	}
}

func NewReturnVariable(chatMessageJSON *[]ChatMessageJSON, text string, code int) *ReturnVariable {
	return &ReturnVariable{
		chatMessageJSON,
		text,
		code,
		nil,
	}
}

func NewReturnVariableString(text string, code int) *ReturnVariable {
	return &ReturnVariable{
		nil,
		text,
		code,
		nil,
	}
}

func NewReturnVariableRoom(room *Room, code int) *ReturnVariable {
	return &ReturnVariable{
		nil,
		"",
		code,
		room,
	}
}