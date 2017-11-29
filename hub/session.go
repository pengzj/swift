package hub

import (
	"io"
	"encoding/base64"
	"crypto/rand"
	"github.com/pengzj/swift/protocol"
	"google.golang.org/grpc"
	"hash/crc32"
	"github.com/pengzj/swift/internal"
	"time"
	"github.com/pengzj/swift/logger"
)

const (
	STATE_UNAVAIABLE = 1
	STATE_AVAIABLE = 2
)


func UniqueId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

type Connection interface {
	Close() error
}

type Session struct {
	Id string
	Conn Connection
	Send chan []byte
	state int
	handlerMap map[string]func()

	userData map[string]interface{}
}

func (session *Session) Bind(name string, handler func())  {
	if session.handlerMap[name] != nil {
		logger.Fatal("func " + name + "bind twice!")
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
		routeName, err := GetHandlerName(routeId)

		var data []byte
		if err != nil {
			data = NotFound()
			session.Write(protocol.Encode(protocol.TYPE_DATA_RESPONSE, protocol.MessageEncode(requestId, routeId, data)))

		} else {
			handler, _ := GetHandler(routeId)
			beforeHandler := GetBeforeHandler()

			if beforeHandler != nil {
				if err = beforeHandler(session, routeName); err != nil {
					session.Write(protocol.Encode(protocol.TYPE_KICK, []byte(err.Error())))
					//wait until the last message  is sent
					time.AfterFunc(50 * time.Millisecond, func() {
						session.Close()
					})
					return
				}
			}

			data = handler(session, in)
			session.Write(protocol.Encode(protocol.TYPE_DATA_RESPONSE, protocol.MessageEncode(requestId, routeId, data)))
		}
	}
}

func (session *Session) Close()  {
	if session.state == STATE_UNAVAIABLE {
		return
	}
	session.state = STATE_UNAVAIABLE
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
	handler := GetRouteHandler(serverType)
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