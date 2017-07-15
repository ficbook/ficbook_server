package main

import (
	"log"
	"net/http"

	"./chat"
)

func main() {
	log.SetFlags(log.Lshortfile)

	// websocket server
	server := chat.NewServer("/entry")
	go server.Listen()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
