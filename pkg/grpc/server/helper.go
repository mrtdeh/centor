package grpc_server

import (
	"context"
	"fmt"
	"time"

	"github.com/mrtdeh/centor/proto"
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
	if _, ok := a.childs[c.Id]; !ok {
		return fmt.Errorf("child %s is not exist", c.Id)
	}
	delete(a.childs, c.Id)
	return nil
}

// ========================== HEALTH CHECK =========================
func connHealthCheck(s *stream, d time.Duration) {
	for {
		if err := connIsFailed(s.conn); err != nil {
			s.err <- err
			return
		}
		_, err := s.proto.Ping(context.Background(), &proto.PingRequest{})
		if err != nil {
			s.err <- err
			return
		}
		time.Sleep(d)
	}
}

func connIsFailed(conn *grpc.ClientConn) error {
	status := conn.GetState()
	if status == connectivity.TransientFailure ||
		// status == connectivity.Idle ||
		status == connectivity.Shutdown {
		return fmt.Errorf("connection is failed with status %s", status)
	}
	return nil
}
