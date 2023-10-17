package grpc_client

import (
	"fmt"

	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type client struct {
	id   string
	addr string
	conn proto.DiscoveryClient
}

func NewClient(name string) *client {
	return &client{id: name}
}

// dial with server
func (c *client) Dial(addr string) error {

	c.addr = addr

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("error in dial : %s", err.Error())
	}
	c.conn = proto.NewDiscoveryClient(conn)

	return nil
}
