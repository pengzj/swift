package connector

import (
	"sync"
	"../hub"
	"./tcp"
	"./udp"
	"./websocket"
	"./option"
)

type Socket interface {
	Start(hub *hub.Hub, host string, port string)
	Read([]byte)
	Write([]byte)
	Close()
	SetOption(*option.ConnectorOption)
}


type Connector struct {
	option *option.ConnectorOption
	socket Socket
}

func (connector *Connector) SetOption(option *option.ConnectorOption)  {
	connector.option = option
}


func (connector *Connector) Start(connType, host, port string)  {
	switch connType {
	case "tcp":
		connector.socket = new(tcp.TcpSocket)
		break;
	case "udp":
		connector.socket = new(udp.UdpSocket)
		break;
	case "websocket":
		connector.socket = new(websocket.WebSocket)
		break;
	}

	if connector.option != nil {
		connector.socket.SetOption(connector.option)
	}

	var hub *hub.Hub = hub.NewHub()
	go hub.Run()

	connector.socket.Start(hub, host, port)
}

func (connector *Connector) Stop()  {
	connector.socket.Close()
}




