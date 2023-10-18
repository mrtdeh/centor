package grpc_client

import (
	"context"
	"fmt"
	"math"

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
	// master election for servers / best election for clients
	var addr string
	var err error
	if cnf.IsServer {
		addr, err = leaderElect(cnf.ServersAddr)
	} else {
		addr, err = bestElect(cnf.ServersAddr)
	}
	if err != nil {
		return err
	}

	fmt.Println("sel addr : ", addr)

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

func bestElect(addrs []string) (string, error) {

	index := -1
	weight := math.MaxInt32
	for i, a := range addrs {

		conn, err := grpc.Dial(a, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return "", fmt.Errorf("error in dial : %s", err.Error())
		}

		c := proto.NewDiscoveryClient(conn)
		res, err := c.GetInfo(context.Background(), &proto.EmptyRequest{})
		if err != nil {
			// log.Println("error in getInfo : ", err.Error())
			fmt.Println("failed to get info from  :", a)
			continue
		}

		if res.Weight < int32(weight) {
			weight = int(res.Weight)
			index = i
			conn.Close()
		}

	}

	if index == -1 {
		return "", fmt.Errorf("server's are not available")
	}

	return addrs[index], nil

}

func leaderElect(addrs []string) (string, error) {
	for _, a := range addrs {
		conn, err := grpc.Dial(a, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return "", fmt.Errorf("error in dial : %s", err.Error())
		}
		c := proto.NewDiscoveryClient(conn)
		res, err := c.GetInfo(context.Background(), &proto.EmptyRequest{})
		if err != nil {
			return "", fmt.Errorf("error in getInfo : %s", err.Error())
		}

		if res.IsMaster {
			conn.Close()
			return a, nil
		}

	}
	return "", fmt.Errorf("leader not found")
}
