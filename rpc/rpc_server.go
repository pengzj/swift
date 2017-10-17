package rpc

import (
	"net"
	"log"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc"
)

type RpcServer struct {
	Server *grpc.Server
}

var std *RpcServer


func NewServer()  *RpcServer{
	return &RpcServer{
		Server: grpc.NewServer(),
	}
}

func GetServer() *RpcServer {
	return std
}

func (rcpServer *RpcServer)Start(host, port string)  {
	listener, err := net.Listen("tcp", host + ":" + port)
	if err != nil {
		log.Fatal(err)
	}

	reflection.Register(std.Server)
	if err := std.Server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

func (rpcServer *RpcServer) Close()  {
	rpcServer.Server.Stop()
}

func init() {
	std = new(RpcServer)
	std.Server = grpc.NewServer()
}