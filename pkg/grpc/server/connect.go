package grpc_server

import (
	"context"
	"fmt"
	"math"

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
