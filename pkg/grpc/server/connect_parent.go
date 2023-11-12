package grpc_server

import (
	"context"
	"fmt"
	"time"

	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	conn, err := grpc.Dial(si.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("error in dial : %s", err.Error())
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

	// join to parent server
	_, err = a.parent.proto.Join(context.Background(), &proto.JoinMessage{
		Id:   a.id,
		Addr: a.addr,
	})
	if err != nil {
		return fmt.Errorf("error in join to server : %s", err.Error())
	}

	// health check conenction for parent server
	go connHealthCheck(&a.parent.stream, time.Second*2)

	// @@@ we can add Sync function to connect bi-directional to server

	fmt.Printf("Connect to server - ID=%s\n", si.Id)
	return <-a.parentErr()
}
