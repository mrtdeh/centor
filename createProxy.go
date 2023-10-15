package main

import (
	"context"

	"github.com/mrtdeh/centor/proto"
)

func (s *Server) CreateProxy(ctx context.Context, req *proto.CreateProxyRequest) (*proto.Close, error) {

	for _, v := range s.clients {
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
