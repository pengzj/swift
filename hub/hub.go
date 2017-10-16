package hub


type Hub struct {
	Sessions map[string]*Session
	Register chan *Session
	Unregister chan *Session
	Broadcast chan []byte
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
			hub.Sessions[session.Id] = session
		case session := <-hub.Unregister:
			delete(hub.Sessions, session.Id)
			close(session.Send)
		case message := <- hub.Broadcast:
			for _, client := range hub.Sessions {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(hub.Sessions, client)
				}
			}
		}
	}
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