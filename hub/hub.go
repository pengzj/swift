package hub

import (
	"sync"
)

type Hub struct {
	Sessions map[string]*Session
	Register chan *Session
	Unregister chan *Session
	Broadcast chan []byte
	mu sync.Mutex

	routeMap map[string]func(*Session)string
}

func NewHub() *Hub {
	return &Hub{
		Sessions: make(map[string]*Session),
		Register: make(chan []*Session),
		Unregister:make(chan []*Session),
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
			close(session.Send)
		case message := <- hub.Broadcast:
			for _, session := range hub.Sessions {
				session.Write(message)
			}
		}
	}
}

func (hub *Hub) Route(serverType string,  handler func(session *Session) string)  {
	hub.routeMap[serverType] = handler
}

func (hub *Hub) GetRouteHandle(serverType string) func(*Session)string  {
	return hub.routeMap[serverType]
}

func (hub *Hub)GetSessionById(id string) *Session {
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
		Register: make(chan []*Session),
		Unregister:make(chan []*Session),
		Broadcast: make(chan []byte),
	}
}