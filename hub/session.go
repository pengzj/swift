package hub

import (
	"net"
	"io"
	"encoding/base64"
	"crypto/rand"
	"log"
	"../protocol"
)


func UniqueId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

type message struct {
	Id int
	Type int
	Body []byte
}

type Session struct {
	Id string
	Hub *Hub
	Conn net.Conn
	Send chan []byte
	Bind func(string, func())
	Trigger func(string)
	handlerMap map[string]func()

	MessageId int

	userData map[string]interface{}
}

func (session *Session) Bind(name string, handler func())  {
	if session.handlerMap[name] != nil {
		log.Fatal("func " + name + "bind twice!")
	}
	session.handlerMap[name] = handler
}

func (session *Session) Trigger(name string)  {
	var handler func() = session.handlerMap[name]
	if handler != nil {
		handler()
	}
}

func (session *Session) HandleData(data []byte)  {
	packageType, body := protocol.Decode(data)
	switch packageType {
	case protocol.PACKAGE_TYPE_HANDSHAKE:
	case protocol.PACKAGE_TYPE_ACK:
	case protocol.PACKAGE_TYPE_HEARTBEAT:
	case protocol.PACKAGE_TYPE_DATA:
	case protocol.PACKAGE_TYPE_ACK:

	}
}

func (session *Session) Close()  {
	session.Trigger("onClosed")
	session.Conn.Close()
	close(session.Send)
}

func (session *Session) Write(data []byte)  {
	session.Send <- data
}