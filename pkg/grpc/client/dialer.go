package grpc_client

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type client struct {
	id   string
	addr string
	conn proto.DiscoveryClient
}

func NewClient() *client {
	return &client{}
}

// dial with server
func (c *client) Dial(addr string) error {

	c.addr = addr
	id := sha256.Sum256([]byte(time.Now().String()))
	c.id = hex.EncodeToString(id[:])[:6]

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("error in dial : %s", err.Error())
	}
	c.conn = proto.NewDiscoveryClient(conn)

	return nil
}

// this is a function that accept message from server
func (c *client) Follow() error {
	req, err := c.conn.Follow(context.Background())
	if err != nil {
		return fmt.Errorf("failed to follow : %s\n", err.Error())
	}

	req.Send(&proto.FollowerRequest{
		Data: &proto.FollowerRequest_JoinMsg{
			JoinMsg: &proto.JoinMessage{
				Id:   c.id,
				Addr: c.addr,
			},
		},
	})

	for {
		res, err := req.Recv()
		if err != nil {
			return fmt.Errorf("error in recivce : %s", err.Error())
		}

		if p := res.GetProxyRequest(); p != nil {
			processProxyRequest(p)
		}

	}
}

func processProxyRequest(p *proto.ProxyRequest) {
	// portNumber++
	// serverPort := fmt.Sprintf("%d", portNumber)

	// service := ServiceProxy{
	// 	// user values
	// 	Name:             "test-service",
	// 	LocalServicePort: p.ProxyPort,
	// 	// protected values
	// 	gRPCServerPort: serverPort,
	// 	msgCh:          make(chan []byte, 1024),
	// }

	// // serviceMap[p.]

	// go func() {
	// 	service.RunProxy()
	// }()

}
