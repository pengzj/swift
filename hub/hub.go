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