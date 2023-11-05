package grpc_server

import (
	"context"
	"fmt"
	"time"

	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
)

type agent struct {
	id       string
	addr     string
	isServer bool
	isLeader bool
	weight   int
	isReady  bool

	brothers []string
	parent   *parent
	childs   map[string]child
}

func (a *agent) checkParent() {
	for {
		_, err := a.parent.proto.Ping(context.Background(), &proto.PingRequest{})
		if err != nil {
			a.parent.stream.err <- fmt.Errorf("ping error parent : %s", err.Error())
		}
		time.Sleep(time.Second * 2)
	}
}

func (a *agent) parentErr() <-chan error {
	return a.parent.stream.err
}

type stream struct {
	conn  *grpc.ClientConn
	proto proto.DiscoveryClient
	err   chan error
}

type parent struct {
	isLeader bool
	stream
}

type child struct {
	Id       string
	Addr     string
	IsServer bool
	stream
}
