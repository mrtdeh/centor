package grpc_server

import (
	"log"

	grpc_client "github.com/mrtdeh/centor/pkg/grpc/client"
)

// server grpc server and register service
func (s *server) ConnectToMaster(name string, addrs []string) error {
	c := grpc_client.Configs{
		Name:                name,
		IsServer:            true,
		ConnectToOnlyMaster: true,
		ServersAddr:         addrs,
	}
	if err := c.Connect(); err != nil {
		log.Fatal(err)
	}
	return nil
}
