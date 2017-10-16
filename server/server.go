package server

import (
	"../connector"
	"../connector/option"
	"../rpc"
)

type Server struct {
	Type string `json:"type"`
	Id string `json:"id"`
	ClientHost string `json:"clienthost,omitempty"`
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
	server.Connector.Start(server.Type, server.ClientHost, server.Port)
}

func (server *Server) startRpcServer()  {
	server.rpcServer = rpc.GetServer()

	if server.IsMaster == true {
		//todo
		rpc.LoadMaster()
	}

	server.rpcServer.Start(server.Host, server.Port)
}

func (server *Server) Stop()  {
	server.rpcServer.Close()
}