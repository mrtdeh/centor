package grpc_server

import (
	"context"
	"fmt"
	"time"

	"github.com/mrtdeh/centor/proto"
)

func (a *agent) ConnectToChild(c *child) error {
	// dial to child listener
	conn, err := grpc_Dial(c.Addr)
	if err != nil {
		return err
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
		err = grpc_ConnectBack(context.Background(), &cc.stream, a.id)
		if err != nil {
			return fmt.Errorf("error in connect back : %s", err.Error())
		}
		// send added node info to leader
		err := a.syncChangeToLeader(NodeInfo{
			Id:       c.Id,
			Address:  c.Addr,
			IsServer: c.IsServer,
			ParentId: a.id,
		}, ChangeActionAdd)
		if err != nil {
			return fmt.Errorf("error in sync change : %s", err.Error())
		}
		// run health check conenction for this child
		go connHealthCheck(&cc.stream, time.Second*2)
	} else {
		return fmt.Errorf("child you want to check not exist")
	}

	return <-c.childErr()
}
