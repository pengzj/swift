package rpc

import (
	"errors"
	"../pb"
	"context"
	"../hub"
	"../internal"
	"log"
)

type Service struct {

}

func (server *Service) InterOnline(ctx context.Context, in *pb.InterOnlineRequest) (*pb.InterOnlineReply, error) {
	reply := &pb.InterOnlineReply{
		Count: uint32(hub.GetHub().Size()),
	}
	return reply, nil
}

func (s *Service) OnlineStatistics(ctx context.Context, in *pb.OnlineRequest) (*pb.OnlineReply, error) {
	serverType := in.Type
	servers := internal.GetServersByType(serverType)
	if len(servers) == 0 {
		return nil, errors.New("server not exist")
	}
	var reply *pb.OnlineReply = &pb.OnlineReply{
		Total:0,
	}

	for _, s := range servers {
		conn := internal.GetClientConnByServerId(s.Id)
		c := pb.NewServiceClient(conn)

		r, err := c.InterOnline(context.Background(), &pb.InterOnlineRequest{})
		if err != nil {
			log.Fatal("inter online ", err)
			//todo replace log
			continue
		}

		reply.Servers = append(reply.Servers, &pb.OnlineReply_Online{Id:s.Id, Count:r.Count,})
		reply.Total = reply.Total + r.Count
	}
	return reply, nil
}

func (s *Service) Offline(ctx context.Context, in *pb.OfflineRequest) (*pb.OfflineReply, error) {
	servers := internal.GetServers()
	if len(servers) == 0 {
		return nil, errors.New("server not exist")
	}
	var reply *pb.OfflineReply = new(pb.OfflineReply)
	for _, s :=range servers {
		conn := internal.GetClientConnByServerId(s.Id)
		c := pb.NewServiceClient(conn)
		_, err := c.InterOffline(context.Background(), &pb.InterOfflineRequest{})
		if err != nil {
			log.Fatal(err)
		}
	}
	return reply, nil
}

func (s *Service) InterOffline(ctx context.Context, in *pb.InterOfflineRequest) (*pb.InterOfflineReply, error)  {
	reply := new(pb.InterOfflineReply)


	hub.GetHub().Stop()

	return reply, nil
}


func loadServer()  {
	pb.RegisterServiceServer(std.Server, &Service{})
}
