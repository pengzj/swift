package rpc

import (
	"errors"
	"log"
	"../pb"
	"context"
	"../hub"
	"../internal"
)

type Service struct {

}

func (server *Service) InterOnline(ctx context.Context, in *pb.InterOnlineRequest) (*pb.InterOnlineReply, error) {
	reply := &pb.InterOnlineReply{
		Count: uint32(hub.GetHub().Size()),
	}
	return reply, nil
}

func LoadServer()  {
	pb.RegisterOnlineServer(std.Server, &Service{})
}

func (s *Service) OnlineStatistics(ctx context.Context, in *pb.OnlineRequest) (*pb.OnlineReply, error) {
	serverType := in.Type
	servers := internal.GetServersByType(serverType)
	if servers == nil {
		return nil, errors.New("server not exist")
	}
	var reply *pb.OnlineReply
	for _, s := range servers {
		conn := internal.GetClientConnByServerId(s.Id)
		c := pb.NewOnlineClient(conn)

		r, err := c.InterOnline(context.Background(), &pb.InterOnlineRequest{})
		if err != nil {
			log.Fatal(err)
		}
		reply.Servers = append(reply.Servers, &pb.OnlineReply_Online{Id:s.Id, Count:r.Count,})
		reply.Total = reply.Total + r.Count
	}
	return reply, nil
}

func LoadMaster()  {
	pb.RegisterOnlineServer(std.Server,&Service{})
}
