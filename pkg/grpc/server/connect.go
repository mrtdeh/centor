package grpc_server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mrtdeh/centor/proto"
)

const (
	StatusDisconnected = "disconnected"
	StatusConnected    = "connected"
	StatusConnecting   = "connecting"
)

func (a *agent) Connect(stream proto.Discovery_ConnectServer) error {
	var done = make(chan bool, 1)
	var joined bool
	var c *child
	var resCh = make(chan *proto.ConnectMessage, 1)
	var errCh = make(chan error, 1)

	// receive message from in sperated goroutine
	go func() {
		for {
			// receive connect message from child server
			res, err := stream.Recv()
			if err != nil {
				errCh <- err
				break
			}
			resCh <- res
		}
	}()

	// receive connect message from channel
	for {
		select {

		// wait for receive message
		case res := <-resCh:
			c = &child{
				agent: agent{
					id:       res.Id,
					dc:       res.DataCenter,
					addr:     res.Addr,
					isServer: res.IsServer,
				},
			}
			// store client connection
			err := addChild(a, c) // add child
			if err != nil {
				return err
			}
			joined = true
			fmt.Printf("Added new client - ID=%s\n", c.id)

			// Dial back to joined server
			go func() {
				err := a.ConnectToChild(c, done)
				if err != nil {
					a.CloseChild(c)
					errCh <- err
				}
			}()

			// Send child status to leader
			go func() {
				// wait for child to connect done
				<-done
				// then, send changes to leader
				err := a.syncChangeToLeader(NodeInfo{
					Id:       c.id,
					Address:  c.addr,
					IsServer: c.isServer,
					ParentId: a.id,
				}, ChangeActionAdd)
				if err != nil {
					errCh <- fmt.Errorf("error in sync change : %s", err.Error())
				}
			}()

		// wait for error
		case err := <-errCh:
			log.Println("conenct error : ", err.Error())
			if joined {
				// leave child from joined server
				leaveChild(a, c)
				// set error to child stream
				c.stream.err <- fmt.Errorf("client disconnected")
				fmt.Printf("Disconnect client - ID=%s\n", c.id)

				// send change for remove client to leader
				err := a.syncChangeToLeader(NodeInfo{
					Id:       c.id,
					Address:  c.addr,
					IsServer: c.isServer,
				}, ChangeActionRemove)
				if err != nil {
					log.Fatalf("error in sync change : %s", err.Error())
				}
			}
			return err

		} // end select
	} // end for
}

func (a *agent) syncChangeToLeader(ni NodeInfo, action int32) error {
	if a.isLeader {
		err := a.applyChange(ni.Id, ni, action)
		if err != nil {
			return fmt.Errorf("error in applyChange : %s", err.Error())
		}
	} else {
		for {
			if a.parent != nil {
				data, err := json.Marshal(ni)
				if err != nil {
					return err
				}
				_, err = a.parent.proto.Change(context.Background(), &proto.ChangeRequest{
					Change: &proto.ChangeRequest_NodesChange{
						NodesChange: &proto.NodesChange{
							Id:     ni.Id,
							Action: action,
							Data:   string(data),
						},
					},
				})
				if err != nil {
					log.Printf("error in syncChangeToLeader : %s", err.Error())
					continue
				}
				// break out of for
				break
			}
			fmt.Println("retry to sync")
			time.Sleep(time.Millisecond * 100)
		}
	}

	return nil
}
func leaveChild(a *agent, c *child) error {
	// delete(a.childs, c.id)
	if _, exist := a.childs[c.id]; exist {
		return fmt.Errorf("this join id is not exist for leaving : %s", c.id)
	}
	a.childs[c.id].status = StatusDisconnected
	a.weight--
	return nil
}
func addChild(a *agent, c *child) error {
	if _, exist := a.childs[c.id]; exist {
		return fmt.Errorf("this join id already exist : %s", c.id)
	}
	// add requested to childs list
	c.status = StatusConnecting
	a.childs[c.id] = c
	a.weight++

	return nil
}
