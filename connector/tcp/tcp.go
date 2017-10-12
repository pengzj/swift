package tcp

import (
	"net"
	"log"
	"../option"
	"../../hub"
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

		go handleConn(conn.(net.TCPConn))

	}
}

func (socket *TcpSocket)Read()  {

}

func (socket *TcpSocket) Write(data []byte)  {

}

func (socket *TcpSocket) Close()  {

}