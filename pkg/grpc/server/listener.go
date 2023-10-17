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
type ServerOptions struct {
	Name    string
	Host    string
	Port    string
	Replica []string
}
type server struct {
	id          string
	addr        string
	isMaster    bool
	connections map[string]connection
}

var s *server

func NewServer(opt ServerOptions) *server {
	// check whether name exist in cluster or not
	name := opt.Name
	addr := fmt.Sprintf("%s:%s", opt.Host, opt.Port)
	s = &server{
		id:          name,
		addr:        addr,
		connections: make(map[string]connection),
	}
	return s
}

// server grpc server and register service
func (s *server) Listen(addr string) error {

	s.addr = addr
	grpcServer := grpc.NewServer()

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("error creating the server %v", err)
	}
	proto.RegisterDiscoveryServer(grpcServer, s)

	fmt.Println("listen an ", addr)

	return grpcServer.Serve(listener)
}

// func IsMaster() bool {
// 	im, err := c.GetInfo(context.Background(), &proto.EmptyRequest{})
// 	if err != nil {
// 		log.Fatal(fmt.Errorf("error check is_master : %s\n", err.Error()))
// 	}

// 	if !im.IsMaster {
// 		return false
// 	}

// 	return true
// }
