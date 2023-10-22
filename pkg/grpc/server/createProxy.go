package grpc_server

import (
	"context"

	"github.com/mrtdeh/centor/proto"
)

func (s *server) CreateProxy(ctx context.Context, req *proto.CreateProxyRequest) (*proto.Close, error) {
	// search service id in services pool and get that target id

	// iterate connections to find target
	for _, v := range s.connections {
		if req.TargetId == v.Id {

			// send proxy request to target
			v.conn.Send(&proto.LeaderResponse{
				Data: &proto.LeaderResponse_ProxyRequest{
					ProxyRequest: &proto.ProxyRequest{
						ProxyPort: req.TargetServicePort,
					},
				},
			})

			break
		}
	}
	return &proto.Close{}, nil
}
