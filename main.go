package main

import (
	"log"
	"net/http"
	"flag"
	"strconv"
	"path/filepath"
	"github.com/ficbook/ficbook_server/chat"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yanzay/cfg"
	"runtime"
)

func main() {
	configPtr := flag.String("config", "config.cfg", "Path to the configuration file")
	buildDB := flag.Bool("database-init", false, "Update the database table")
	createRoom := flag.String("create-room", "", "Creates a room")
	flag.Parse()

	path, _ := filepath.Abs(*configPtr)
	cfgInfo := make(map[string]string)
	err := cfg.Load(path, cfgInfo)
	if err != nil {
		log.Fatal(err)
	}

	if runtime.GOOS == "windows" {
		path, _ = filepath.Abs(filepath.Dir("locals/" + cfgInfo["server_lang"] + ".cfg"))
	} else {
		path, _ = filepath.Abs(filepath.Dir("locals\\" + cfgInfo["server_lang"] + ".cfg"))
	}
	lang := make(map[string]string)
	err = cfg.Load(path, lang)
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(cfgInfo["db_server"], cfgInfo["db_user"] + ":" + cfgInfo["db_password"] + "@/" + cfgInfo["db_db"] + "?charset=utf8mb4&parseTime=true")
	defer db.Close()

	// websocket server
	server := chat.NewServer(cfgInfo["server_pattern"], db, buildDB, createRoom, &lang)
	go server.Listen()

	// Updating count of users
	updateTime, _ := strconv.Atoi(cfgInfo["update_time"])
	go server.UpdateOnlineRooms(updateTime)
	
	log.Fatal(http.ListenAndServe(cfgInfo["server_ip"] + ":" + cfgInfo["server_port"], nil))
}
