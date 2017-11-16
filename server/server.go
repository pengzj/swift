package server

import (
	"github.com/pengzj/swift/connector"
	"github.com/pengzj/swift/connector/option"
	"github.com/pengzj/swift/rpc"
)

type Server struct {
	Type string `json:"type"`
	ConnType string `json:"connType"`
	Id string `json:"id"`
	ClientHost string `json:"clientHost,omitempty"`
	ClientPort string `json:"clientPort,omitempty"`
	Host string `json:"host"`
	Port string `json:"port"`
	Frontend bool `json:"frontend"`

	IsMaster bool `json:"-"`
	Connector *connector.Connector `json:"-"`
	rpcServer *rpc.RpcServer `json:"-"`
}

func (server *Server) Start(option *option.ConnectorOption)  {
	if server.Frontend == true {
		server.startServer(option)
	}

	server.startRpcServer()
}

func (server *Server) startServer(option *option.ConnectorOption)  {
	server.Connector = new(connector.Connector)
	server.Connector.SetOption(option)
	server.Connector.Start(server.ConnType, server.ClientHost, server.ClientPort)
}

func (server *Server) startRpcServer()  {
	server.rpcServer = rpc.GetServer()

	server.rpcServer.Start(server.Host, server.Port)
}

func (server *Server) Stop()  {
	if server.Frontend == true {
		server.Connector.Close()
	}
	server.rpcServer.Close()
}