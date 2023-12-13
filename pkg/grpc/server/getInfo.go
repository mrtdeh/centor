package grpc_server

import (
	"context"

	"github.com/mrtdeh/centor/proto"
)

func (s *Agent) GetInfo(context.Context, *proto.EmptyRequest) (*proto.InfoResponse, error) {
	return &proto.InfoResponse{
		Id:       s.id,
		IsLeader: s.isLeader,
		Weight:   int32(s.weight),
	}, nil
}
