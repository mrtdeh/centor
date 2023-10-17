package grpc_client

import (
	"context"
	"fmt"

	"github.com/mrtdeh/centor/proto"
)

// this is a function that accept message from server
func (c *client) follow() error {
	req, err := c.conn.Follow(context.Background())
	if err != nil {
		return fmt.Errorf("failed to follow : %s\n", err.Error())
	}

	req.Send(&proto.FollowerRequest{
		Data: &proto.FollowerRequest_JoinMsg{
			JoinMsg: &proto.JoinMessage{
				Id:       c.id,
				Addr:     c.serverAddr,
				IsServer: c.isServer,
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
	fmt.Println("proxy request from server : ", p.ProxyPort)
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
