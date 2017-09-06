package server

import (
	"math/rand"
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"encoding/json"
	"strconv"
)

func (s *Server) Rooms_List(w http.ResponseWriter, r *http.Request) {
	var rooms []Room
	for _, room := range(s.rooms) {
		rooms = append(rooms, *room)
	}
	bytes, _ := json.Marshal(&rooms)
	w.Write(bytes)
}

func (s *Server) Rooms_GetRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if room, ok := s.rooms[id]; ok && err == nil {
		b, _ := json.Marshal(room)
		w.Write(b)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("This room does not exist"))
	}
}

func (s *Server) Users_SignIn(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	m := make(map[string]interface{})
	m["id"] = rand.Intn(10)
	m["token"] = RandStringRunes(20)
	m["username"] = r.Form.Get("login")
	m["level"] = 0
	b, _ := json.Marshal(&m)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func GetName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Print(vars["name"])
}
