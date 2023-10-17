package grpc_client

import (
	"fmt"

	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Configs struct {
	Name                string
	IsServer            bool
	ConnectToOnlyMaster bool
	ServersAddr         []string
}
type client struct {
	id         string
	serverAddr string
	isServer   bool
	conn       proto.DiscoveryClient
}

// connect to server
func (cnf *Configs) Connect() error {
	// master election
	addr := cnf.ServersAddr[0]
	// dial to selected master
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("error in dial : %s", err.Error())
	}
	// create client object
	c := &client{
		id:         cnf.Name,
		serverAddr: addr,
		isServer:   cnf.IsServer,
		conn:       proto.NewDiscoveryClient(conn),
	}
	// client do follow master
	if err := c.follow(); err != nil {
		return fmt.Errorf("error in follow : %s", err.Error())
	}

	return nil
}
