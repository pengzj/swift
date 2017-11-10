package connector

import (
	"../hub"
	"./tcp"
	"./websocket"
	"./option"
)

type Socket interface {
	Start(host string, port string)
	Close()
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
		connector.socket = &tcp.TcpSocket{
			CloseChan:make(chan bool),
		}
	case "ws":
		connector.socket = &websocket.WebSocket{
			CloseChan:make(chan bool),
			PongChan:make(chan bool),
		}
	}


	go hub.GetHub().Run()

	go connector.socket.Start(host, port)
}

func (connector *Connector) Close()  {
	connector.socket.Close()
}




