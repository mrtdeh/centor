package grpc_server

import (
	"fmt"
	"net"

	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
)

type connection struct {
	Id   string
	Addr string
	conn proto.Discovery_FollowServer
}
type server struct {
	id          string
	isMaster    bool
	connections map[string]connection
}

var s *server

// server grpc server and register service
func Serve(addr string, register ...func(*grpc.Server)) error {

	grpcServer := grpc.NewServer()

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("error creating the server %v", err)
	}
	proto.RegisterDiscoveryServer(grpcServer, s)
	// register(grpcServer)

	return grpcServer.Serve(listener)
}
