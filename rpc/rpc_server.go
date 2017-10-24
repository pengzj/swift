package rpc

import (
	"net"
	"log"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc"
	"../internal"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/codes"
	"fmt"
	"context"
)

type RpcServer struct {
	Server *grpc.Server
}

var std *RpcServer

var rpcHandler func()


func GetServer() *RpcServer {
	return std
}

func auth(ctx context.Context) error {
	md , ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return grpc.Errorf(codes.Unauthenticated, "no token")
	}

	var token string
	if val, ok := md["token"]; ok {
		token = val[0]
	}
	fmt.Println(md)
	if token != internal.GetSecretKey() {
		return grpc.Errorf(codes.Unauthenticated, "token invalid")
	}

	return nil
}


func (rcpServer *RpcServer)Start(host, port string)  {
	listener, err := net.Listen("tcp", host + ":" + port)
	if err != nil {
		log.Fatal(err)
	}

	var interceptor grpc.UnaryServerInterceptor = func(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		err = auth(ctx)
		if err != nil {
			return
		}

		return handler(ctx, req)
	}

	var opts []grpc.ServerOption
	opts = append(opts, grpc.UnaryInterceptor(interceptor))

	std.Server = grpc.NewServer(opts...)

	loadServer()

	if rpcHandler != nil {
		rpcHandler()
	}

	reflection.Register(std.Server)
	if err := std.Server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}



func RegisterRPC(handler func())  {
	if rpcHandler != nil {
		panic("rpc handle has register twice!")
	}
	rpcHandler = handler
}

func (rpcServer *RpcServer) Close()  {
	rpcServer.Server.Stop()
}

func init() {
	std = new(RpcServer)
}