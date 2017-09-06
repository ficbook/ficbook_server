package server

import (
	"net/http"
 	"net/url"
 	"io/ioutil"
	"encoding/json"
)

type isLogin struct {
	result bool
	data string
}

//Authorization is authorized on the site ficbook.net. Returns true if the login and password are correct.
func Authorization(login string, password string) bool {
	resp, _ := http.PostForm("https://ficbook.net/login_check", url.Values{"login": {login}, "password": {password}})
	body, _ := ioutil.ReadAll(resp.Body)
	var dat map[string]interface{}
	json.Unmarshal(body, &dat)
	return dat["result"].(bool)
}

func GetIDFicbookProfile(login string) uint32 {
	resp, _ := http.PostForm("https://ficbook.net/ajax/user_info", url.Values{"nickname": {login}})
	body, _ := ioutil.ReadAll(resp.Body)
	var dat map[string]interface{}
	json.Unmarshal(body, &dat)
	if dat["result"].(bool) == true {
		data := dat["data"].(map[string]interface{})
		return uint32(data["id"].(float64))
	}
	return 0
}