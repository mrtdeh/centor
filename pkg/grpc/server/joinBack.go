package grpc_server

import (
	"context"
	"fmt"

	"github.com/mrtdeh/centor/proto"
)

func (s *agent) JoinBack(ctx context.Context, req *proto.JoinBackMessage) (*proto.Close, error) {
	fmt.Printf("server join back - ID=%s\n", req.Id)

	return &proto.Close{}, nil
}
