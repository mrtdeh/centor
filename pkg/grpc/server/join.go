package grpc_server

import (
	"context"
	"fmt"

	"github.com/mrtdeh/centor/proto"
)

func (a *agent) Join(ctx context.Context, req *proto.JoinMessage) (*proto.Close, error) {
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
			fmt.Println(err.Error())
			a.CloseChild(c)
		}
	}()

	return &proto.Close{}, nil
}
