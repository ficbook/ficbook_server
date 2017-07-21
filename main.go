package main

import (
	"log"
	"net/http"
	"flag"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	//"github.com/ficbook/ficbook_server/chat"
	"./chat_d"
	"github.com/yanzay/cfg"
)

func main() {
	log.SetFlags(log.Lshortfile)

	configPtr := flag.String("config", "config.cfg", "Path to the configuration file")
	flag.Parse()

	cfgInfo := make(map[string]string)
	err := cfg.Load(*configPtr, cfgInfo)
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(cfgInfo["db_server"], cfgInfo["db_user"] + ":" + cfgInfo["db_password"] + "@/" + cfgInfo["db_table"])
	defer db.Close()

	// websocket server
	server := chat.NewServer(cfgInfo["server_pattern"], db)
	go server.Listen()

	log.Fatal(http.ListenAndServe(cfgInfo["server_ip"] + ":" + cfgInfo["server_port"], nil))
}
