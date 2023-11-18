package grpc_server

import (
	"fmt"

	"github.com/mrtdeh/centor/proto"
)

func (a *agent) ConnectBack(stream proto.Discovery_ConnectBackServer) error {

	for {
		res, err := stream.Recv()
		if err != nil {
			return err
		}

		pid := res.Id

		defer func() {
			a.parent.stream.err <- fmt.Errorf("parent disconnected")
			fmt.Printf("Disconnect parent - ID=%s\n", pid)
		}()

		fmt.Printf("Conenct Back from parent - ID=%s\n", pid)
	}
}
