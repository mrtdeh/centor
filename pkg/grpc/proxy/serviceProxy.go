package grpc_proxy

import (
	"log"
	"net"

	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
)

type Configs struct {
	Name             string
	LocalServicePort string
	GRPCServerPort   string
}

type server struct {
	msgIn  chan []byte
	msgOut chan []byte
	alive  bool
}

func (cnf *Configs) Listen() {

	s := &server{alive: true}

	addr := "localhost:" + cnf.GRPCServerPort
	grpcServer := grpc.NewServer()

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("error creating the server %v", err)
	}
	proto.RegisterProxyManagerServer(grpcServer, s)

	go s.tcpDialToService(cnf.LocalServicePort)

	if err := grpcServer.Serve(listener); err != nil {
		s.alive = false
		log.Println("error in client proxy server : ", err.Error())
	}

}

func (s *server) tcpDialToService(port string) {

}
