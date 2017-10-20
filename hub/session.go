package hub

import (
	"net"
	"io"
	"encoding/base64"
	"crypto/rand"
	"log"
	"../protocol"
	"google.golang.org/grpc"
	"hash/crc32"
	"../internal"
	"fmt"
	"time"
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
	fmt.Println(packageType, string(body), time.Now())
	switch packageType {
	case protocol.TYPE_HANDSHAKE:
		data = GetHandlers()
		message := protocol.MessageEncode(0, 0, data)
		session.Write(protocol.Encode(protocol.TYPE_HANDSHAKE_ACK, message))
		return
	case protocol.TYPE_DATA_NOTIFY:
		_, routeId, body := protocol.MessageDecode(body)
		handler, err := GetHandler(routeId)
		if err == nil {
			_ = handler(session, body)
		}
	case protocol.TYPE_DATA_REQUEST:
		requestId, routeId, in := protocol.MessageDecode(body)
		handler, err := GetHandler(routeId)
		var data []byte
		fmt.Println("recieve request ", time.Now())
		if err != nil {
			data = NotFound()
		} else {
			data = handler(session, in)
		}
		fmt.Println("recieve resposne ", time.Now())
		session.Write(protocol.Encode(protocol.TYPE_DATA_RESPONSE, protocol.MessageEncode(requestId, routeId, data)))
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
	handlerId, err := GetHandlerId(route)
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
	handler := GetRouteHandle(serverType)
	var serverId string
	if handler == nil {
		servers :=  internal.GetServersByType(serverType)
		crc := crc32.ChecksumIEEE([]byte(session.Id))
		idx := int(crc%uint32(len(servers)))
		serverId = servers[idx].Id
	} else {
		serverId = handler(session)
	}
	return internal.GetClientConnByServerId(serverId)
}