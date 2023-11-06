package grpc_server

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

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

// ======================================
func (a *agent) checkParent() {
	for {
		status := a.parent.conn.GetState()
		if status != connectivity.Ready {
			a.parent.stream.err <- fmt.Errorf("error parent connection status : %s", status)
			return
		}
		time.Sleep(time.Second * 2)
	}
}

func (a *agent) parentErr() <-chan error {
	return a.parent.stream.err
}

// =======================================
func (c *child) checkChild() {
	for {
		status := c.conn.GetState()
		// fmt.Println("status : ", status)
		if status == connectivity.TransientFailure ||
			status == connectivity.Idle ||
			status == connectivity.Shutdown {
			c.stream.err <- fmt.Errorf("error child connection status : %s", status)
			return
		}
		time.Sleep(time.Second * 2)
	}
}

func (c *child) childErr() <-chan error {
	return c.stream.err
}

// ======================================
func (a *agent) waitForReady() {
	for {
		if a.isReady {
			return
		}
		time.Sleep(time.Second)
	}
}

func (a *agent) ready() {
	a.isReady = true
}

func (a *agent) unReady() {
	a.isReady = false
}

func (a *agent) CloseChild(c *child) error {
	// return c.conn.Close()
	child, ok := a.childs[c.Id]
	if !ok {
		return fmt.Errorf("child %s is not exist", child.Id)
	}
	delete(a.childs, child.Id)
	return nil
}
