package connector

import (
	"../hub"
	"./tcp"
	"./udp"
	"./websocket"
	"./option"
)

type Socket interface {
	Start(host string, port string)
	Write([]byte)
	Close()
	SetOption(*option.ConnectorOption)
	readPump()
	writePump()
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

	go hub.GetHub().Run()

	connector.socket.Start(host, port)
}

func (connector *Connector) Close()  {
	connector.socket.Close()
}




