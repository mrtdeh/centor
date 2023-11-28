package grpc_server

import (
	"fmt"
	"log"

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
					isLeader: res.IsLeader,
				},
			}
			// store client connection
			err := addChild(a, c) // add child
			if err != nil {
				errCh <- err
			}
			joined = true

			// send back agent id to connected client
			err = stream.Send(&proto.ConnectBackMessage{Id: a.id})
			if err != nil {
				errCh <- err
			}

			// Dial back to joined server
			go func() {
				err := a.CreateChildStream(c, done)
				if err != nil {
					errCh <- err
				}
			}()

			// Send child status to leader
			go func() {
				// wait for child to connect done
				<-done
				// then, send changes to leader
				err := a.syncAgentChange(&c.agent, ChangeActionAdd)
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

				// send change for remove client to leader
				err := a.syncAgentChange(&c.agent, ChangeActionRemove)
				if err != nil {
					log.Fatalf("error in sync change : %s", err.Error())
				}
			}
			return err

		} // end select
	} // end for
}

func leaveChild(a *agent, c *child) {
	if _, exist := a.childs[c.id]; !exist {
		fmt.Printf("this join id is not exist for leaving : %s", c.id)
		return
	}
	a.childs[c.id].status = StatusDisconnected
	a.weight--
	c.stream.err <- fmt.Errorf("client disconnected")
	fmt.Printf("Disconnect client - ID=%s\n", c.id)
}
func addChild(a *agent, c *child) error {
	if cc, exist := a.childs[c.id]; exist && cc.status == StatusConnected {
		return fmt.Errorf("this join id already exist : %s", c.id)
	}
	// add requested to childs list
	c.status = StatusConnecting
	a.childs[c.id] = c
	a.weight++

	fmt.Printf("Added new client - ID=%s\n", c.id)
	return nil
}
