package grpc_server

import (
	"fmt"
	"log"

	"github.com/mrtdeh/centor/proto"
)

func (s *server) Follow(stream proto.Discovery_FollowServer) error {
	var client connection
	res, err := stream.Recv()
	if err != nil {
		log.Fatal(err)
	}

	if j := res.GetJoinMsg(); j != nil {
		if _, ok := s.connections[j.Id]; ok {
			return fmt.Errorf("name is select by another nodes")
		}

		client := connection{
			conn:   stream,
			Id:     j.Id,
			Addr:   j.Addr,
			Server: j.IsServer,
		}

		s.connections[j.Id] = client
		s.weight++

		fmt.Printf("client added ID=%s IsServer=%v \n", j.Id, j.IsServer)
	}

	return <-client.err
}
