package grpc_server

import (
	"context"

	"github.com/mrtdeh/centor/proto"
)

func (s *server) GetInfo(context.Context, *proto.EmptyRequest) (*proto.InfoResponse, error) {
	return &proto.InfoResponse{
		Id:       s.id,
		IsMaster: s.isMaster,
	}, nil
}
