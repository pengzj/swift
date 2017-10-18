package hub

import (
	"net"
	"io"
	"encoding/base64"
	"crypto/rand"
	"log"
	"../protocol"
	"../internal"
	"google.golang.org/grpc"
	"hash/crc32"
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
	Conn net.Conn
	Send chan []byte
	handlerMap map[string]func()

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
	case protocol.TYPE_HANDSHAKE:
		data = internal.GetRoutes()
		message := protocol.MessageEncode(0, 0, data)
		session.Write(protocol.Encode(protocol.TYPE_HANDSHAKE_ACK, message))
		return
	case protocol.TYPE_HEARTBEAT:
		return
	case protocol.TYPE_DATA_NOTIFY:
		_, routeId, body := protocol.MessageDecode(body)
		handler, err := internal.GetHandler(routeId)
		if err == nil {
			_ = handler(session, body)
		}
		return
	case protocol.TYPE_DATA_REQUEST:
		requestId, routeId, body := protocol.MessageDecode(body)
		handler, err := internal.GetHandler(routeId)
		var data []byte
		if err != nil {
			data = internal.NotFound()
		} else {
			data = handler(session, body)
		}
		session.Write(protocol.Encode(protocol.TYPE_DATA_RESPONSE, protocol.MessageEncode(requestId, routeId, data)))
		return
	}
}

func (session *Session) Close()  {
	GetHub().Unregister <- session
	session.Trigger("onClosed")
	session.Conn.Close()
	close(session.Send)
}

func (session *Session) Write(data []byte)  {
	session.Send <- data
}

func (session *Session) Push(route string, data []byte)  {
	handlerId, err := internal.GetHandlerId(route)
	if err != nil {
		return
	}
	session.Write(protocol.Encode(protocol.TYPE_DATA_PUSH, protocol.MessageEncode(0, handlerId, data)))
}

func (session *Session) Kick(data []byte)  {
	session.Write(protocol.Encode(protocol.TYPE_KICK, protocol.MessageEncode(0,0, data)))
}

func (session *Session) Set(key string, value interface{})  {
	session.userData[key] = value
}

func (session *Session) Get(key string) interface{} {
	return session.userData[key]
}

func (session *Session) GetClientConn(serverType string) *grpc.ClientConn  {
	handler := GetHub().GetRouteHandle(serverType)
	var serverId string
	if handler == nil {
		servers := internal.GetServersByType(serverType)
		crc := crc32.ChecksumIEEE([]byte(session.Id))
		idx := int(crc%uint32(len(servers)))
		serverId = servers[idx].Id

	} else {
		serverId = handler(session)
	}
	return internal.GetClientConnByServerId(serverId)
}