package udp

import (
	"github.com/pengzj/swift/connector/option"
)

//todo udp will not be supported in a while, we will focus tcp/websocket

type UdpSocket struct {

}

func (socket *UdpSocket)Start( host string, port string)  {

}


func (socket *UdpSocket) SetOption(option *option.ConnectorOption)  {

}

func (socket *UdpSocket) Close()  {

}