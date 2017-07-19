package main

import (
	"log"
	"net/http"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/ficbook/ficbook_server/chat"
)

func main() {
	log.SetFlags(log.Lshortfile)

	db, _ := gorm.Open("mysql", "root:qwer99@/ficbook?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()

	// websocket server
	server := chat.NewServer("/", db)
	go server.Listen()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
