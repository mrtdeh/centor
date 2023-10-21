package grpc_client

import (
	"context"
	"fmt"

	grpc_proxy "github.com/mrtdeh/centor/pkg/grpc/proxy"
	"github.com/mrtdeh/centor/proto"
)

var (
	done           = make(chan bool, 1)
	portNumber int = 11000
)

// this is a function that accept message from server
func (c *client) follow() error {
	req, err := c.conn.Follow(context.Background())
	if err != nil {
		return fmt.Errorf("failed to follow : %s\n", err.Error())
	}
	// send join message to server
	err = req.Send(&proto.FollowerRequest{
		Data: &proto.FollowerRequest_JoinMsg{
			JoinMsg: &proto.JoinMessage{
				Id: c.id,
				// Addr:     c.,
				IsServer: c.isServer,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send : %s\n", err.Error())
	}

	// set joining done
	done <- true

	// now ready for recieve any request from server
	for {
		res, err := req.Recv()
		if err != nil {
			return fmt.Errorf("error in recivce : %s", err.Error())
		}

		if p := res.GetProxyRequest(); p != nil {
			go processProxyRequest(p)
		}

	}
}

func (c *client) wait() {
	<-done
}

func processProxyRequest(p *proto.ProxyRequest) {
	fmt.Println("proxy request from server : ", p.ProxyPort)
	portNumber++
	serverPort := fmt.Sprintf("%d", portNumber)

	service := grpc_proxy.Configs{
		// Name:             "test-service",
		LocalServicePort: p.ProxyPort,
		GRPCServerPort:   serverPort,
	}

	// serviceMap[p.]

	service.Listen()

}
