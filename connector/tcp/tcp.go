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
	heartbeatInterval = 5 * time.Second
)

type TcpSocket struct {
	hub.Session
	Conn *net.TCPConn
	option *option.ConnectorOption
}

func (socket *TcpSocket) Start(h *hub.Hub, host ,port string)  {
	tcpAddr, err := net.ResolveTCPAddr("tcp", host + ":" + port)
	if err != nil {
		log.Fatal(err)
	}
	listener, err := net.Listen("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		conn = conn.(net.TCPConn)
		session := &TcpSocket{hub.Session{
			Id: hub.UniqueId(),
			Hub: h,
			Conn:&conn,
			Send: make(chan []byte),
		}}
		session.Hub.Register <- session

		go session.readPump()
		go session.writePump()

	}
}

func (socket *TcpSocket)readPump()  {
	var buffer bytes.Buffer
	var length int
	for {
		data := make([]byte, 64)
		_, err := socket.Conn.Read(data)
		if err != nil {
			if err == net.UnknownNetworkError.Error() {
				return
			}
		}
		buffer.Write(data)
		length = protocol.GetPackageLength(buffer.Bytes())
		message := make([]byte, length)
		_, err = buffer.Read(message)
		if err != nil {
			log.Fatalf("read data error: %v", err)
		}
		socket.HandleData(message)
	}
}

func (socket *TcpSocket) writePump()  {
	ticker := time.NewTicker(heartbeatInterval)
	defer func() {
		ticker.Stop()
		socket.Close()
	}()

	for {
		select {
		case message, ok := socket.Send:
			socket.Conn.SetWriteDeadline(time.Now().Add(heartbeatInterval))
			if !ok {
				//todo close session
				return
			}

			_, err := socket.Conn.Write(message)
			if err != nil {
				//todo close
				return
			}

		case <-ticker.C:
			socket.Conn.SetWriteDeadline(time.Now().Add(heartbeatInterval))
			//todo write heartbeat
		}
	}
}


func (socket *TcpSocket) SetOption(option *option.ConnectorOption)  {
	socket.option = option
}