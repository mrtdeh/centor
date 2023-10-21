package grpc_server

import (
	"context"
	"log"

	"github.com/mrtdeh/centor/proto"
)

func (s *server) CreateProxy(ctx context.Context, req *proto.CreateProxyRequest) (*proto.Close, error) {
	log.Println("debug 1 : ", s.connections)
	for _, v := range s.connections {
		log.Println("debug 1.5 : ", v.Id, req.TargetId)
		if req.TargetId == v.Id {
			log.Println("debug 2")
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
