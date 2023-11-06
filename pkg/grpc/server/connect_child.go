package grpc_server

import (
	"fmt"
	"time"

	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
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
			c := a.childs[c.Id]
			for {
				status := c.conn.GetState()
				if status == connectivity.TransientFailure ||
					status == connectivity.Idle ||
					status == connectivity.Shutdown {
					c.stream.err <- fmt.Errorf("error child connection status : %s", status)
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
