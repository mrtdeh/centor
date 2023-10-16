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
	err  chan error
	conn proto.Discovery_FollowServer
}
type server struct {
	id          string
	isMaster    bool
	connections map[string]connection
}

var s *server

func NewServer(name string) *server {

	s = &server{
		name,
		true,
		make(map[string]connection),
	}
	return s
}

// server grpc server and register service
func (s *server) Serve(addr string) error {

	grpcServer := grpc.NewServer()

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("error creating the server %v", err)
	}
	proto.RegisterDiscoveryServer(grpcServer, s)

	fmt.Println("listen an ", addr)

	return grpcServer.Serve(listener)
}
