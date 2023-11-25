package grpc_server

import (
	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
)

type agent struct {
	id        string // id of the agent
	addr      string // address of this node
	isServer  bool   // is this node a server or not
	isLeader  bool   // is this node leader or not
	isPrimary bool   // is this node primary server or not
	isReady   bool   // is this node ready or not
	weight    int    // weight of this node in the cluster

	servers []string          // servers in the cluster
	parent  *parent           // parent of this node in the cluster or in primary cluster
	childs  map[string]*child // childs of this node in the cluster
}

type stream struct {
	conn  *grpc.ClientConn      // connection to the server
	proto proto.DiscoveryClient // discovery protocol
	err   chan error            // channel for error
	close chan bool             // channel for closed connection
}

type parent struct {
	id        string // id of the parent
	isLeader  bool   // is this node leader or not
	isPrimary bool   // is this node primary server or not
	stream
}

type child struct {
	Id       string // id of the child
	Addr     string // address of the child
	IsServer bool   // is this node a server or not
	stream          // stream of the child
}
