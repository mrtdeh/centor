package grpc_server

import (
	"context"
	"fmt"
	"time"

	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
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
		si, err = leaderElect(a.servers)
	} else {
		si, err = bestElect(a.servers)
	}
	if err != nil {
		return err
	}

	// dial to selected master
	conn, err := grpc.Dial(si.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("error in dial : %s", err.Error())
	}
	defer conn.Close()

	// create parent object
	a.parent = &parent{
		stream: stream{
			conn:  conn,
			proto: proto.NewDiscoveryClient(conn),
			err:   make(chan error, 1),
		},
	}

	// call join rpc to parent server
	_, err = a.parent.proto.Join(context.Background(), &proto.JoinMessage{
		Id:   a.id,
		Addr: a.addr,
	})
	if err != nil {
		return fmt.Errorf("error in join to server : %s", err.Error())
	}

	// run health check service for parent server
	go func() {
		for {
			status := a.parent.conn.GetState()
			if status != connectivity.Ready {
				a.parent.stream.err <- fmt.Errorf("error parent connection status : %s", status)
				return
			}
			time.Sleep(time.Second * 2)
		}
	}()

	fmt.Printf("Connect to server - ID=%s\n", si.Id)
	return <-a.parentErr()
}