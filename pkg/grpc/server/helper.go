package grpc_server

import (
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

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
func ConnIsFailed(conn *grpc.ClientConn) error {
	status := conn.GetState()
	if status == connectivity.TransientFailure ||
		status == connectivity.Idle ||
		status == connectivity.Shutdown {
		return fmt.Errorf("connection is failed with status %s", status)
	}
	return nil
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
