package main

import (
	"context"
	"log"
	"net"

	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
)

type ServiceProxy struct {
	msgCh            chan []byte
	alive            bool
	Name             string
	LocalServicePort string
	gRPCServerPort   string
}

type ServerProxy struct{}

var serverProxy *ServerProxy

// client read payload bytes from other client
func (sp *ServerProxy) SendPayload(ctx context.Context, req *proto.RequestPayload) (*proto.Close, error) {
	// req.Conn
	return nil, nil
}

func (s *ServiceProxy) RunProxy() {

	addr := "localhost:" + s.gRPCServerPort
	grpcServer := grpc.NewServer()

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("error creating the server %v", err)
	}

	portNumber++
	s.alive = true

	proto.RegisterProxyManagerServer(grpcServer, serverProxy)
	if err := grpcServer.Serve(listener); err != nil {
		s.alive = false
		log.Println("error in client proxy server : ", err.Error())
	}

}
