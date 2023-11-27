package grpc_server

import (
	"context"

	"github.com/mrtdeh/centor/proto"
)

func (s *agent) Exec(context.Context, *proto.ExecRequest) (*proto.ExecResponse, error) {
	return &proto.ExecResponse{}, nil
}
