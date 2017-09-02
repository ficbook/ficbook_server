package server

//APIReturn - the structure returned by the ParseAPI function
type APIReturn struct {
	Interface *map[string]interface{}
}

func NewAPIReturn(Interface *map[string]interface{}) *APIReturn {
	return &APIReturn{
		Interface,
	}
}