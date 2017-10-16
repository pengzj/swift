package rpc

import (
	"context"
	"../pb"
	"../hub"
)

type Server struct {

}

func (server *Server) InterOnline(ctx context.Context, in *pb.InterOnlineRequest) (*pb.InterOnlineReply, error) {
	reply := &pb.InterOnlineReply{
		Count: hub.GetHub().Size(),
	}
	return reply, nil
}

func LoadServer()  {
	pb.RegisterOnlineServer(std.Server, &Server{})
}