package grpc_server

import (
	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
)

type Agent struct {
	id        string // id of the agent
	addr      string // address of this node
	dc        string // datacenter of this node
	isServer  bool   // is this node a server or not
	isLeader  bool   // is this node leader or not
	isPrimary bool   // is this node primary server or not
	isReady   bool   // is this node ready or not
	weight    int    // weight of this node in the cluster

	parent *Parent           // parent of this node in the cluster or in primary cluster
	childs map[string]*Child // childs of this node in the cluster
}

type stream struct {
	conn  *grpc.ClientConn      // connection to the server
	proto proto.DiscoveryClient // discovery protocol
	err   chan error            // channel for error
	close chan bool             // channel for closed connection
}

type Parent struct {
	Agent  // parent agent information
	stream // stream of parent server
}

type Child struct {
	Agent         // child agent information
	stream        // stream of the child server
	status string // status of child in the cluster
}

// ===========================================
