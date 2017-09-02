package main

import (
	"github.com/gorilla/websocket"
	"github.com/ficbook/ficbook_server/newserver/chat"
	"net/http"
	"log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yanzay/cfg"
)

var (
	updater = websocket.Upgrader{}
)

func main() {

	cfgInfo := make(map[string]string)
	err := cfg.Load("config.cfg", cfgInfo)
	if err != nil {
		log.Fatal(err)
	}

	db, _ := gorm.Open(cfgInfo["db_server"], cfgInfo["db_user"] + ":" + cfgInfo["db_password"] + "@/" + cfgInfo["db_db"] + "?charset=utf8mb4&parseTime=true")
	defer db.Close()

	s := server.NewServer(db)
	go s.Listen()
	log.Fatal(http.ListenAndServe(cfgInfo["server_ip"] + ":" + cfgInfo["server_port"], nil))
}