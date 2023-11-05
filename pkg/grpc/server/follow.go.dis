package grpc_server

import (
	"context"
	"fmt"

	grpc_proxy "github.com/mrtdeh/centor/pkg/grpc/proxy"
	"github.com/mrtdeh/centor/proto"
)

var (
	portNumber int = 11000
)

func (s *agent) Follow(stream proto.Discovery_FollowServer) error {
	var c child
	for {
		res, err := stream.Recv()
		if err != nil {
			fmt.Printf("client closed ID=%s\n", c.Id)
			delete(s.childs, c.Id)
			s.weight--
			return err
		}

		if j := res.GetJoinMsg(); j != nil {
			if _, ok := s.childs[j.Id]; ok {
				return fmt.Errorf("name is select by another nodes")
			}

			c = child{
				conn:     stream,
				Id:       j.Id,
				Addr:     j.Addr,
				IsServer: j.IsServer,
			}

			s.childs[j.Id] = c
			s.weight++

			fmt.Printf("client added ID=%s IsServer=%v \n", j.Id, j.IsServer)
		}
	}
}

// this is a function that accept message from server
func (a *agent) follow() error {
	req, err := a.parent.conn.Follow(context.Background())
	if err != nil {
		return fmt.Errorf("failed to follow : %s\n", err.Error())
	}
	// send join message to server
	err = req.Send(&proto.FollowerRequest{
		Data: &proto.FollowerRequest_JoinMsg{
			JoinMsg: &proto.JoinMessage{
				Id: a.id,
				// Addr:     a.,
				IsServer: a.isServer,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send : %s\n", err.Error())
	}

	// set ready to true
	a.ready()
	defer a.unReady()

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

func processProxyRequest(p *proto.ProxyRequest) {
	fmt.Println("proxy request from server : ", p.ProxyPort)
	portNumber++
	serverPort := fmt.Sprintf("%d", portNumber)

	service := grpc_proxy.Configs{
		LocalServicePort: p.ProxyPort,
		GRPCServerPort:   serverPort,
	}

	service.Listen()

}
