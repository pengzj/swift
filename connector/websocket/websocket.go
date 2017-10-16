package websocket

import (
	"github.com/gorilla/websocket"
	"./../../hub"
	"net/http"
	"log"
	"time"
	"bytes"
	"../../protocol"
	"../option"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:1024,
	WriteBufferSize:1024,
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

type WebSocket struct {
	hub.Session
	Conn *websocket.Conn
	option *option.ConnectorOption
	certFile string
	keyFile string
}

func (socket *WebSocket) SetCert(certFile, keyFile string)  {
	socket.certFile = certFile
	socket.keyFile = keyFile
}

func (socket *WebSocket)Start(host string, port string)  {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r)
	})

	var err error
	if socket.certFile != nil && socket.keyFile != nil {
		err = http.ListenAndServeTLS(host + ":" + port, socket.certFile, socket.keyFile, nil)
	} else {
		err = http.ListenAndServe(host + ":" + port, nil)
	}
	if err != nil {
		log.Fatal(err)
	}

}

func serveWs(w http.ResponseWriter, r *http.Request)  {
	conn, err := upgrader.Upgrade(w,r, nil)
	if err != nil {
		log.Fatal(err)
	}

	session := &WebSocket{hub.Session{
		Id: hub.UniqueId(),
		Conn:conn,
		Send:make(chan []byte),
	}}
	hub.GetHub().Register <- session

	go session.readPump()
	go session.writePump()
}


func (socket *WebSocket) readPump()  {
	defer func() {
		socket.Close()
	}()
	socket.Conn.SetReadLimit(maxMessageSize)
	socket.Conn.SetReadDeadline(time.Now().Add(pongWait))
	socket.Conn.SetPongHandler(func(string) error {
		socket.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})


	var buffer bytes.Buffer
	var length int
	for {
		_, message, err := socket.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Fatalf("error: %v", err)
			}
			break
		}

		buffer.Write(message)

		length = protocol.GetPackageLength(buffer.Bytes())
		data := make([]byte, length)
		_, err = buffer.Read(data)
		if err != nil {
			log.Fatalf("read data error: %v", err)
		}
		socket.HandleData(data)
	}
}

func (socket *WebSocket) writePump()  {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		socket.Close()
	}()
	for {
		select {
		case message, ok := socket.Send:
			socket.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				socket.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := socket.Conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			socket.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := socket.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Fatal(err)
				return
			}

		}
	}
}

func (socket *WebSocket) SetOption(option *option.ConnectorOption)  {
	socket.option = option
}