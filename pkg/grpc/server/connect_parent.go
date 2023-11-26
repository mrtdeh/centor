package grpc_server

import (
	"context"
	"fmt"
	"time"

	"github.com/mrtdeh/centor/proto"
)

type connectConfig struct {
	ServersAddresses []string
	ConnectToPrimary bool
}

func (a *agent) ConnectToParent(cc connectConfig) error {
	if len(cc.ServersAddresses) == 0 {
		return nil
	}
	servers := cc.ServersAddresses

	// master election for servers / best election for clients
	var si *ServerInfo
	var err error
	if a.isServer {
		// select leader only
		si, err = leaderElect(servers)
	} else {
		// select best server in server's pool
		si, err = bestElect(servers)
	}
	if err != nil {
		return err
	}

	// dial to selected server
	conn, err := grpc_Dial(si.Addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// create parent object
	a.parent = &parent{
		agent: agent{ // parent agent
			id:        si.Id,
			isLeader:  si.IsLeader,
			isPrimary: cc.ConnectToPrimary,
		},
		stream: stream{ // parent stream
			conn:  conn,
			proto: proto.NewDiscoveryClient(conn),
			err:   make(chan error, 1),
		},
	}
	if a.isLeader {
		if n, ok := nodesInfo[a.id]; ok {
			n.ParentId = si.Id
			nodesInfo[a.id] = n
		}
	}

	// create sync stream rpc to parent server
	err = grpc_Connect(context.Background(), a)
	if err != nil {
		return fmt.Errorf("error in sync : %s", err.Error())
	}

	// health check conenction for parent server
	go connHealthCheck(&a.parent.stream, time.Second*2)

	// @@@ we can add Sync function to connect bi-directional to server

	return <-a.parentErr()
}
