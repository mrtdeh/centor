package grpc_server

import (
	"context"
	"fmt"
	"time"

	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (a *agent) ConnectToChild(c *child) error {
	// dial to child listener
	conn, err := grpc.Dial(c.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("error in dial : %s", err.Error())
	}
	defer conn.Close()

	// create child object
	if cc, ok := a.childs[c.Id]; ok {
		cc.stream = stream{
			conn:  conn,
			proto: proto.NewDiscoveryClient(conn),
			err:   make(chan error, 1),
		}
		a.childs[c.Id] = cc

		go func() {
			client := a.childs[c.Id]
			for {
				if err := ConnIsFailed(client.conn); err != nil {
					client.stream.err <- fmt.Errorf("Closed client - ID=%s", client.Id)
					return
				}
				_, err := client.proto.Ping(context.Background(), &proto.PingRequest{})
				if err != nil {
					a.parent.stream.err <- fmt.Errorf("error parent ping : %s", err.Error())
					return
				}
				time.Sleep(time.Second * 2)
			}
		}()
	} else {
		return fmt.Errorf("child you want to check not exist")
	}

	// run health check service for this child

	fmt.Printf("Added new client - ID=%s\n", c.Id)
	return <-c.childErr()
}
