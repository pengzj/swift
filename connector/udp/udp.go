package udp

import (
	"../../hub"
)

//todo udp will not be supported in a while, we will focus tcp/websocket

type UdpSocket struct {

}

func (socket *UdpSocket)Start(hub *hub.Hub, host string, port string)  {

}

func (socket *UdpSocket)Read()  {

}

func (socket *UdpSocket) Write(data []byte)  {

}

func (socket *UdpSocket) Close()  {

}