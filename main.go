package main

import (
	"log"
	"net"

	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
)


type Server struct {
	Id       string
	IsMaster bool
	clients  map[string]Client
}
type Client struct {
	Id   string
	Addr string
	conn proto.Discovery_FollowServer
}

var server *Server

func main() {
	addr := "localhost:90001"
	grpcServer := grpc.NewServer()

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("error creating the server %v", err)
	}

	proto.RegisterDiscoveryServer(grpcServer, server)
	grpcServer.Serve(listener)
}
