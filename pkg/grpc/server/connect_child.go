package grpc_server

import (
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
		// store client connection and proto info
		cc.stream = stream{
			conn:  conn,
			proto: proto.NewDiscoveryClient(conn),
			err:   make(chan error, 1),
		}
		// run health check conenction for this child
		go connHealthCheck(&cc.stream, time.Second*2)
	} else {
		return fmt.Errorf("child you want to check not exist")
	}

	fmt.Printf("Added new client - ID=%s\n", c.Id)
	return <-c.childErr()
}
