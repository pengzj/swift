package rpc

import (
	"net"
	"log"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc"
)

type RpcServer struct {
	server *grpc.Server
}

var std *RpcServer


func GetServer()  *RpcServer{
	if std == nil {
		std = new(RpcServer)
		std.server = grpc.NewServer()
	}
	return std
}

func (rcpServer *RpcServer)Start(host, port string)  {
	listener, err := net.Listen("tcp", host + ":" + port)
	if err != nil {
		log.Fatal(err)
	}

	//todo register server handle

	reflection.Register(std)
	if err := std.server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

func (rcpServer *RpcServer) Close()  {
	
}