package grpc_server

import (
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
