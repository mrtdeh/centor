package grpc_server

import (
	"fmt"

	"github.com/mrtdeh/centor/proto"
)

func (a *agent) Sync(stream proto.Discovery_SyncServer) error {
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
				fmt.Printf("Disconnect client - ID=%s\n", c.Id)
				a.weight--
			}
		}()

		// check for getting Join message
		if join := res.GetJoinMsg(); join != nil {
			// create new child from join request
			c = &child{
				Id:       join.Id,
				Addr:     join.Addr,
				IsServer: join.IsServer,
			}
			// check weather requested id is exist or not
			if _, exist := a.childs[join.Id]; exist {
				return fmt.Errorf("this join id already exist : %s", join.Id)
			}
			// add requested to childs list
			a.childs[join.Id] = c
			a.weight++
			joined = true

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
}
