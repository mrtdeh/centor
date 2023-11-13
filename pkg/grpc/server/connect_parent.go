package grpc_server

import (
	"context"
	"fmt"
	"time"

	"github.com/mrtdeh/centor/proto"
)

func (a *agent) ConnectToParent() error {
	if len(a.servers) == 0 {
		return nil
	}

	// master election for servers / best election for clients
	var si *ServerInfo
	var err error
	if a.isServer {
		// select leader only
		si, err = leaderElect(a.servers)
	} else {
		// select best server in server's pool
		si, err = bestElect(a.servers)
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
		id: si.Id,
		stream: stream{
			conn:  conn,
			proto: proto.NewDiscoveryClient(conn),
			err:   make(chan error, 1),
		},
	}

	// create sync stream rpc to parent server
	err = grpc_SyncToParent(context.Background(), &a.parent.stream, a.id, a.addr)
	if err != nil {
		return fmt.Errorf("error in sync : %s", err.Error())
	}

	// health check conenction for parent server
	go connHealthCheck(&a.parent.stream, time.Second*2)

	// @@@ we can add Sync function to connect bi-directional to server

	fmt.Printf("Connect to server - ID=%s\n", si.Id)
	return <-a.parentErr()
}
