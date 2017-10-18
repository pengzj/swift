package tcp

import (
	"net"
	"log"
	"../option"
	"../../hub"
	"bytes"
	"../../protocol"
	"time"
)

var (
	heartbeatInterval = 10 * time.Second
)

type TcpSocket struct {
	option *option.ConnectorOption
	switcher chan bool
}

func (socket *TcpSocket) Start(host ,port string)  {
	tcpAddr, err := net.ResolveTCPAddr("tcp", host + ":" + port)
	if err != nil {
		log.Fatal(err)
	}
	listener, err := net.Listen("tcp", tcpAddr.String())
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		select {
		case <- socket.switcher:
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				log.Fatal(err)
			}

			session := &hub.Session{
				Id: hub.UniqueId(),
				Conn: conn,
				Send: make(chan []byte),
			}

			hub.GetHub().Register <- session

			go readPump(session)
			go writePump(session)
		}

	}
}

func readPump(session *hub.Session)  {
	var buffer bytes.Buffer
	var headerLength = protocol.GetHeadLength()
	var currentTotalLength int
	var length int
	for {
		data := make([]byte, 1024)
		n, err := session.Conn.Read(data)
		if err != nil {
			log.Fatal("read data error: %v", err)
			return
		}
		buf := make([]byte, n)
		copy(buf, data[0:n])

		buffer.Write(buf)

		currentTotalLength = len(buffer.Bytes())
		length = headerLength +  protocol.GetBodyLength(buffer.Bytes())
		message := make([]byte, length)
		if length <= currentTotalLength {
			_, err = buffer.Read(message)
			if err != nil {
				log.Fatalf("read data error: %v", err)
			}
			session.HandleData(message)

			leftLength := currentTotalLength - length
			if leftLength > 0 {
				leftData := make([]byte, leftLength)
				_, err = buffer.Read(leftData)
				if err != nil {
					log.Fatal("package data error: %v", err)
				}
				buffer.Reset()
				buffer.Write(leftData)
			}
		}
	}
}

func writePump(session *hub.Session)  {
	ticker := time.NewTicker(heartbeatInterval)
	defer func() {
		ticker.Stop()
		session.Close()
	}()

	for {
		select {
		case message := <-session.Send:
			session.Conn.SetWriteDeadline(time.Now().Add(heartbeatInterval))

			_, err := session.Conn.Write(message)
			if err != nil {
   				return
			}

		case <-ticker.C:
			session.Conn.SetWriteDeadline(time.Now().Add(heartbeatInterval))
			_, err := session.Conn.Write(protocol.Encode(protocol.TYPE_HEARTBEAT, []byte{}))
			if err != nil {
				log.Fatal(err)
				return
			}
		}
	}
}


func (socket *TcpSocket) SetOption(option *option.ConnectorOption)  {
	socket.option = option
}

func (socket *TcpSocket) Close()  {
	socket.switcher <- false
}