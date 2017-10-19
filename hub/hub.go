package hub

import (
	"sync"
	"encoding/json"
)

type Hub struct {
	Sessions map[string]*Session
	Register chan *Session
	Unregister chan *Session
	Broadcast chan []byte
	mu sync.Mutex

}

func NewHub() *Hub {
	return &Hub{
		Sessions: make(map[string]*Session),
		Register: make(chan *Session),
		Unregister:make(chan *Session),
		Broadcast: make(chan []byte),
	}
}

var std *Hub

func (hub *Hub) Run()  {
	for {
		select {
		case session := <-hub.Register:
			hub.mu.Lock()
			hub.Sessions[session.Id] = session
			hub.mu.Unlock()
		case session := <-hub.Unregister:
			hub.mu.Lock()
			delete(hub.Sessions, session.Id)
			hub.mu.Unlock()
		case message := <- hub.Broadcast:
			for _, session := range hub.Sessions {
				session.Write(message)
			}
		}
	}
}

func (hub *Hub)GetSessionById(id string) *Session {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	return hub.Sessions[id]
}

func (hub *Hub) Kick(session *Session)  {
	hub.Unregister <- session
}

func (hub *Hub) Size() int {
	return len(hub.Sessions)
}

func GetHub()  *Hub {
	return std
}

func init() {
	std = &Hub{
		Sessions: make(map[string]*Session),
		Register: make(chan *Session),
		Unregister:make(chan *Session),
		Broadcast: make(chan []byte),
	}

	inter = new(interEntity)
	inter.handlerList = []string{}
	inter.handlerMap = map[string]func(*Session, []byte) []byte {}
	inter.routeMap = map[string]func(*Session)string{}


	handler := func(*Session, []byte) []byte {
		var body = struct {
			Code int
			Message string
		}{
			Code: 404,
			Message:"method not exists",
		}
		data, _ := json.Marshal(body)
		return data
	}

	HandleFunc("notFound", handler)
}