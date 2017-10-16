package rpc

import (
	"context"
	"../pb"
	"../hub"
)

func (service *Service) InterOnline(ctx context.Context, in *pb.InterOnlineRequest) (*pb.InterOnlineReply, error) {
	reply := &pb.InterOnlineReply{
		Count: hub.GetHub().Size(),
	}
	return reply, nil
}
