package hub

import (
	"sync"
	"encoding/json"
)

type Hub struct {
	Sessions map[string]*Session
	Register chan *Session
	Unregister chan *Session
	Closed chan bool
	Count int
	mu sync.Mutex

}

func NewHub() *Hub {
	return &Hub{
		Sessions: make(map[string]*Session),
		Register: make(chan *Session),
		Unregister:make(chan *Session),
		Closed: make(chan bool),
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
			hub.Count++
		case session := <-hub.Unregister:
			hub.mu.Lock()
			delete(hub.Sessions, session.Id)
			hub.mu.Unlock()
			hub.Count--
		case  <- hub.Closed:
			return
		}
	}
}

func (hub *Hub) Stop()  {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	for _, session :=range hub.Sessions {
		session.Close()
		delete(hub.Sessions, session.Id)
	}


	hub.Closed <- true

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
	return hub.Count
}

func GetHub()  *Hub {
	return std
}

func init() {
	std = NewHub()

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