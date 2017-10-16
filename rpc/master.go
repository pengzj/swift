package rpc

import (
	"../pb"
	"context"
	"../internal"
	"log"
)

type Master struct {

}


func (s *Master) OnlineStatistics(ctx context.Context, in *pb.OnlineRequest) (*pb.OnlineReply, error) {
	serverType := in.Type
	servers := internal.GetServersByType(serverType)
	if servers == nil {
		return nil, error.Error("server not exist")
	}
	var reply *pb.OnlineReply
	for _, server := range servers {
		conn := internal.GetClientConnByServer(server)
		c := pb.NewOnlineClient(conn)

		r, err := c.InterOnline(context.Background(), &pb.InterOnlineRequest{})
		if err != nil {
			log.Fatal(err)
		}
		reply.Servers = append(reply.Servers, &pb.OnlineReply_Online{Id:server.Id, Count:r.Count,})
		reply.Total = reply.Total + r.Count
	}
	return reply, nil
}

func LoadMaster()  {
	pb.RegisterOnlineServer(std,&Master{})
}
