package grpc_server

import (
	"log"
	"net"

	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
)

func NewServer(addr string) {

	grpcServer := grpc.NewServer()

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("error creating the server %v", err)
	}

	proto.RegisterDiscoveryServer(grpcServer, server)
	grpcServer.Serve(listener)
}
