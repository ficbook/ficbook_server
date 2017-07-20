package main

import (
	"log"
	"net/http"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/ficbook/ficbook_server/chat"
	"github.com/yanzay/cfg"
)

func main() {
	log.SetFlags(log.Lshortfile)

	cfgInfo := make(map[string]string)
	err := cfg.Load("config.cfg", cfgInfo)
	if err != nil {
		log.Fatal(err)
	}

	db, _ := gorm.Open(cfgInfo["db_server"], cfgInfo["db_user"] + ":" + cfgInfo["db_password"] + "@/" + cfgInfo["db_table"] + "?charset=" + cfgInfo["charset"] + "&parseTime=" + cfgInfo["parse_time"] + "&loc=" + cfgInfo["loc"])
	defer db.Close()

	// websocket server
	server := chat.NewServer(cfgInfo["server_pattern"], db)
	go server.Listen()

	log.Fatal(http.ListenAndServe(cfgInfo["server_ip"] + ":" + cfgInfo["server_port"], nil))
}
