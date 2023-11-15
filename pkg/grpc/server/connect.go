package grpc_server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mrtdeh/centor/proto"
)

func (a *agent) Connect(stream proto.Discovery_ConnectServer) error {
	var joined bool
	var c *child
	var resCh = make(chan *proto.ConnectMessage, 1)
	var errCh = make(chan error, 1)

	go func() {
		for {
			res, err := stream.Recv()
			if err != nil {
				errCh <- err
				break
			}
			resCh <- res
		}
	}()

	for {
		select {

		// wait for receive message
		case res := <-resCh:
			c = &child{
				Id:       res.Id,
				Addr:     res.Addr,
				IsServer: res.IsServer,
			}
			// check weather requested id is exist or not
			if _, exist := a.childs[res.Id]; exist {
				return fmt.Errorf("this join id already exist : %s", res.Id)
			}
			// add requested to childs list
			a.childs[res.Id] = c
			a.weight++
			joined = true
			fmt.Printf("Added new client - ID=%s\n", c.Id)

			// Dial back to joined server
			go func() {
				err := a.ConnectToChild(c)
				if err != nil {
					a.CloseChild(c)
					errCh <- err
				}
			}()

		// wait for error
		case err := <-errCh:
			if joined {
				a.weight--
				c.stream.err <- fmt.Errorf("client disconnected")
				fmt.Printf("Disconnect client - ID=%s\n", c.Id)

				// send change for remove client to leader
				err := a.syncChangeToLeader(NodeInfo{
					Id:       c.Id,
					Address:  c.Addr,
					IsServer: c.IsServer,
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
		if a.parent != nil {
			data, err := json.Marshal(ni)
			if err != nil {
				return err
			}
			_, err = a.parent.proto.Change(context.Background(), &proto.ChangeRequest{
				Change: &proto.ChangeRequest_NodesChange{
					NodesChange: &proto.NodesChange{
						Id:     ni.Id,
						Action: ChangeActionAdd,
						Data:   string(data),
					},
				},
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
