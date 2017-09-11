package server

import (
	"strings"
	"math/rand"
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"encoding/json"
	"strconv"
	"time"
	"sort"
)

func (s *Server) Rooms_List(w http.ResponseWriter, r *http.Request) {
	var (
		rooms []Room
		keys []int
	)
	for k := range(s.rooms) {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, key := range(keys) {
		rooms = append(rooms, *s.rooms[key])
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

func (s *Server) Rooms_Create(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if len(r.Form.Get("name")) > 0 {
		room := Room{
			Name: r.Form.Get("name"),
			CreatedAt: time.Now(),
		}
		s.db.Create(&room)
		s.db.First(&room).Order("id desc")
		fmt.Println(room)
		s.rooms[room.ID] = &room
		b, _ := json.Marshal(&room)
		w.Write(b)
	}
}

func (s *Server) Rooms_Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idRoom, err := strconv.Atoi(vars["id"])
	if room, ok := s.rooms[idRoom]; ok && err == nil {
		s.db.Delete(room)
		delete(s.rooms, idRoom)
	}
}

func (s *Server) Messages_Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idRoom, err := strconv.Atoi(vars["id"])
	if _, ok := s.rooms[idRoom]; ok && err == nil {
		var messages []ChatMessage
		s.db.Find(&messages).Where("room_id = ?", idRoom).Order("id desc").Count(20)
		messagesJSON, _ := json.Marshal(&messages)
		w.Write(messagesJSON)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("This room does not exist"))
	}
}

func (s *Server) Messages_Post(w http.ResponseWriter, r *http.Request) {
	token := w.Header().Get("X-Auth-Token")
	if len(token) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("The token is missing in the header"))
	} else {
		vars := mux.Vars(r)
		idRoom, err := strconv.Atoi(vars["id"])
		if _, ok := s.rooms[idRoom]; ok && err == nil {
			r.ParseForm()
			body := r.Form.Get("body")
			if !strings.Contains(body, "") {
				user := User{
					Token: token,
				}
				s.db.First(&user)
				if !user.CreatedAt.IsZero() {
					message := ChatMessage{
						UserID: user.ID,
						RoomID: idRoom,
						Body: body,
						CreatedAt: time.Now(),
					}
					s.db.Create(&message)
				} else {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte("The user does not exist with this token"))	
				}				
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("This room does not exist"))
		}
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
