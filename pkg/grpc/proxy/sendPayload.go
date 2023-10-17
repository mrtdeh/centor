package grpc_proxy

import (
	"context"

	"github.com/mrtdeh/centor/proto"
)

// client read payload bytes from other client
func (sp *server) SendPayload(ctx context.Context, req *proto.RequestPayload) (*proto.Close, error) {
	// req.Conn

	return nil, nil
}
