package websocket

import (
	"github.com/gorilla/websocket"
	"time"
	"net/http"
	"log"
	"github.com/pengzj/swift/hub"
	"github.com/pengzj/swift/protocol"
	"bytes"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:1024,
	WriteBufferSize:1024,
}

var (
	heartbeatInterval = 10 * time.Second
)

type WebSocket struct {
	server *http.Server
	CloseChan chan bool
	PongChan chan bool
	certFile string
	keyFile string
}

const (
	//time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	//Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

func (socket *WebSocket) Start(host, port string)  {
	socket.server = &http.Server{Addr: host + ":" + port, Handler:nil}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r)
	})

	var err error
	if len(socket.certFile) > 0 && len(socket.keyFile) > 0 {
		err = socket.server.ListenAndServeTLS(socket.certFile, socket.keyFile)
	} else {
		err = socket.server.ListenAndServe()
	}
	if err != nil {
		log.Fatal(err)
	}
}

func (socket *WebSocket) Close()  {
	close(socket.CloseChan)
	socket.server.Close()
}

func serveWs(w http.ResponseWriter, r *http.Request)  {
	conn, err := upgrader.Upgrade(w,r, nil)
	if err != nil {
		log.Fatal(err)
	}

	session := &hub.Session{
		Id: hub.UniqueId(),
		Conn: conn,
		Send: make(chan []byte, 10),
	}
	hub.GetHub().Register <- session

	go readDump(session)
	go writeDump(session)
}

func readDump(session *hub.Session)  {
	conn := session.Conn.(*websocket.Conn)
	conn.SetReadLimit(int64(1024 * 1024 * 1024))
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			return
		}

		switch messageType {
		case websocket.BinaryMessage, websocket.TextMessage:
			session.HandleData(p)
		case websocket.PingMessage:
			err = conn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(pongWait))
			if err != nil {
				return
			}
		case websocket.CloseMessage:
			session.Close()
			return
		}
	}
}

func writeDump(session *hub.Session)  {
	conn := session.Conn.(*websocket.Conn)
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()
	var buffer bytes.Buffer
	for {
		select {
		case message := <-session.Send:
			buffer.Write(message)
			n := len(session.Send)

			for i := 0; i < n; i++ {
				buffer.Write(<-session.Send)
			}

			conn.SetWriteDeadline(time.Now().Add(heartbeatInterval))
			err := conn.WriteMessage(websocket.BinaryMessage, buffer.Bytes())
			if err != nil {
				return
			}
			buffer.Reset()
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(heartbeatInterval))
			err := conn.WriteMessage(websocket.BinaryMessage, protocol.Encode(protocol.TYPE_HEARTBEAT, []byte{}))
			if err != nil {
				return
			}
		}
	}
}

