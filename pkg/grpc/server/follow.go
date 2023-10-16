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
		client := connection{
			conn: stream,
			Id:   j.Id,
			Addr: j.Addr,
		}

		s.connections[j.Id] = client

		fmt.Println("client added : ", j.Id)
	}

	return <-client.err
}
