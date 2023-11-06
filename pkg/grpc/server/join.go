package grpc_server

import (
	"context"
	"log"

	"github.com/mrtdeh/centor/proto"
)

func (a *agent) Join(ctx context.Context, req *proto.JoinMessage) (*proto.JoinResponse, error) {
	// Dial back to joined server
	c := &child{
		Id:       req.Id,
		Addr:     req.Addr,
		IsServer: req.IsServer,
	}
	a.childs[req.Id] = c
	// Inreament weight of server
	a.weight++

	go func() {
		err := a.ConnectToChild(c)
		if err != nil {
			log.Printf("child %s is dead : %s\n", c.Id, err.Error())
			a.CloseChild(c)
		}
	}()

	res := &proto.JoinResponse{ServerId: a.id}

	return res, nil
}
