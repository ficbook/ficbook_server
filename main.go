package main

import (
	"log"
	"net/http"
	"flag"
	"strconv"
	"os"
	"io"
	"path/filepath"
	"github.com/ficbook/ficbook_server/chat"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yanzay/cfg"
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

	locals, err := http.Get(cfgInfo["path_lang"])
	if err != nil {
		log.Fatal(err)
	}
	defer locals.Body.Close()
	fileLang, err := os.Create("tmplocfb")
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(fileLang, locals.Body)
	if err != nil {
		log.Fatal(err)
	}
	lang := make(map[string]string)
	err = cfg.Load("tmplocfb", lang)
	if err != nil {
		log.Fatal(err)
	}
	fileLang.Close()
	os.Remove("tmplocfb")

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
