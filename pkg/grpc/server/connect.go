package grpc_server

import (
	"context"
	"fmt"

	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (a *agent) ConnectToParent() error {
	// master election for servers / best election for clients
	var addr string
	var err error
	if a.isServer {
		addr, err = leaderElect(a.brothers)
	} else {
		addr, err = bestElect(a.brothers)
	}
	if err != nil {
		return err
	}

	// dial to selected master
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("error in dial : %s", err.Error())
	}
	// create client object
	a.parent = &parent{
		stream: stream{
			conn:  conn,
			proto: proto.NewDiscoveryClient(conn),
			err:   make(chan error, 1),
		},
	}

	_, err = a.parent.proto.Join(context.Background(), &proto.JoinMessage{
		Id:   a.id,
		Addr: a.addr,
	})
	if err != nil {
		return fmt.Errorf("error in join to server : %s", err.Error())
	}

	go a.checkParent()

	return <-a.parentErr()
}
