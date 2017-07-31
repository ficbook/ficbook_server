package main

import (
	"log"
	"fmt"
	"net/http"
	"flag"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/ficbook/ficbook_server/chat"
	//"./chat_d"
	"github.com/yanzay/cfg"
)

func parseCommand(s *string) {
	for true {
		fmt.Print(">>> ")
		fmt.Scanln(s)
		fmt.Println(*s)
	}
}

func main() {
	log.SetFlags(log.Lshortfile)

	configPtr := flag.String("config", "config.cfg", "Path to the configuration file")
	buildDB := flag.Bool("database-init", false, "Update the database table")
	createRoom := flag.String("create-room", "", "Creates a room")
	flag.Parse()

	cfgInfo := make(map[string]string)
	err := cfg.Load(*configPtr, cfgInfo)
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(cfgInfo["db_server"], cfgInfo["db_user"] + ":" + cfgInfo["db_password"] + "@/" + cfgInfo["db_db"] + "?charset=utf8mb4&parseTime=true")
	defer db.Close()

	// websocket server
	server := chat.NewServer(cfgInfo["server_pattern"], db, buildDB, createRoom)
	go server.Listen()

	var stringCommand string
	go parseCommand(&stringCommand)

	log.Fatal(http.ListenAndServe(cfgInfo["server_ip"] + ":" + cfgInfo["server_port"], nil))
}
