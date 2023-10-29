package grpc_server

import (
	"fmt"
	"log"
	"net"

	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
)

type Config struct {
	Name     string
	Host     string
	Port     uint
	Replica  []string
	IsServer bool
	IsLeader bool
}
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

func Init(cnf Config) error {

	s := &agent{
		id:       cnf.Name,
		addr:     fmt.Sprintf("%s:%d", cnf.Host, cnf.Port),
		childs:   make(map[string]child),
		done:     make(chan bool, 1),
		isServer: cnf.IsServer,
		isLeader: cnf.IsLeader,
		brothers: cnf.Replica,
	}

	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("error creating the server %v", err)
	}
	proto.RegisterDiscoveryServer(grpcServer, s)
	fmt.Println("listen an ", s.addr)

	if !cnf.IsLeader && len(cnf.Replica) > 0 {
		go func() {
			if err := s.Connect(); err != nil {
				log.Fatal(err)
			}

		}()
	}

	return grpcServer.Serve(listener)
}
