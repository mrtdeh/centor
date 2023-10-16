package grpc_server

import (
	"context"

	"github.com/mrtdeh/centor/proto"
)

func (s *server) CreateProxy(ctx context.Context, req *proto.CreateProxyRequest) (*proto.Close, error) {

	for _, v := range s.connections {
		if req.TargetAddr == v.Addr {

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
	return nil, nil
}
