package grpc_server

import (
	"fmt"

	"github.com/mrtdeh/centor/proto"
)

func (a *agent) Connect(stream proto.Discovery_ConnectServer) error {
	var joined bool
	var c *child
	for {
		res, err := stream.Recv()
		if err != nil {
			return err
		}
		defer func() {
			// decrease weight when client connection have error
			if joined {
				a.weight--
				c.stream.err <- fmt.Errorf("client disconnected")
				fmt.Printf("Disconnect client - ID=%s\n", c.Id)
			}
		}()

		// create new child from join request
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
				fmt.Println(err.Error())
				a.CloseChild(c)
			}
		}()

	}
}
