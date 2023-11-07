package grpc_server

import (
	"context"
	"log"

	"github.com/mrtdeh/centor/proto"
)

func (a *agent) Call(ctx context.Context, req *proto.CallRequest) (*proto.CallResponse, error) {
	var tags []string
	tags = append(tags, a.id)

	if a.parent != nil && a.parent.id != req.AgentId {

		res, err := a.parent.proto.Call(context.Background(), &proto.CallRequest{
			AgentId: a.id,
		})
		if err != nil {
			log.Fatal("error in call parent :", err.Error())
		}
		tags = append(tags, res.Tags...)
	}

	if a.childs != nil {
		for _, c := range a.childs {
			if c.Id != req.AgentId {
				res, err := c.proto.Call(context.Background(), &proto.CallRequest{
					AgentId: a.id,
				})
				if err != nil {
					log.Fatalf("error in call child %s : %s", c.Id, err.Error())
				}
				tags = append(tags, res.Tags...)
			}
		}
	}

	return &proto.CallResponse{Tags: tags}, nil
}
