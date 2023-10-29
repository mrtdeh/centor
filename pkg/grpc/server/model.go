package grpc_server

import "github.com/mrtdeh/centor/proto"

type agent struct {
	id       string
	addr     string
	isServer bool
	isLeader bool
	weight   int
	parent   *parent // connection client -> server
	done     chan bool

	brothers []string
	childs   map[string]child
}
type parent struct {
	conn     proto.DiscoveryClient
	isLeader bool
}

type child struct {
	Id       string
	Addr     string
	IsServer bool
	err      chan error
	conn     proto.Discovery_FollowServer
}
