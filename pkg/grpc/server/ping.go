package grpc_server

import (
	"context"

	"github.com/mrtdeh/centor/proto"
)

func (s *Agent) Ping(context.Context, *proto.PingRequest) (*proto.PongResponse, error) {
	return &proto.PongResponse{}, nil
}
