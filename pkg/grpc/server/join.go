package grpc_server

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func (s *agent) Join(ctx context.Context, req *proto.JoinMessage) (*proto.Close, error) {
	conn, err := grpc.Dial(req.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &proto.Close{}, fmt.Errorf("error in dial : %s", err.Error())
	}

	c := child{
		stream: stream{
			conn:  conn,
			proto: proto.NewDiscoveryClient(conn),
		},
		Id:       req.Id,
		Addr:     req.Addr,
		IsServer: req.IsServer,
	}

	go func() {
		for {
			_, err := c.proto.JoinBack(context.Background(), &proto.JoinBackMessage{
				Id: s.id,
			})
			if err != nil {
				s := status.Convert(err)
				log.Println("error in joinBack : ", c.Addr, s.Code())
			}
			time.Sleep(time.Second * 10)
		}
	}()

	s.childs[req.Id] = c
	fmt.Printf("Added new client - ID=%s\n", req.Id)
	s.weight++
	return &proto.Close{}, nil
}
