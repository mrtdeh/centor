package grpc_server

import (
	"fmt"
	"net"

	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
)

type connection struct {
	Id     string
	Addr   string
	Server bool
	err    chan error
	conn   proto.Discovery_FollowServer
}
type Configs struct {
	Name     string
	Host     string
	Port     uint
	Replica  []string
	IsLeader bool
}
type server struct {
	id          string
	addr        string
	isMaster    bool
	weight      int
	master      proto.DiscoveryClient
	connections map[string]connection
}

// server grpc server and register service
func (cnf *Configs) Listen() error {

	s := &server{
		addr:        fmt.Sprintf("%s:%d", cnf.Host, cnf.Port),
		connections: make(map[string]connection),
		isMaster:    cnf.IsLeader,
	}

	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("error creating the server %v", err)
	}
	proto.RegisterDiscoveryServer(grpcServer, s)
	fmt.Println("listen an ", s.addr)

	if len(cnf.Replica) > 0 {
		go s.ConnectToMaster(cnf.Name, cnf.Replica)
	}

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
