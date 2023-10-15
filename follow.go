package main

import (
	"log"

	"github.com/mrtdeh/centor/proto"
)

func (s *Server) Follow(stream proto.Discovery_FollowServer) error {

	res, err := stream.Recv()
	if err != nil {
		log.Fatal(err)
	}

	if j := res.GetJoinMsg(); j != nil {
		client := Client{
			conn: stream,
			Id:   j.Id,
			Addr: j.Addr,
		}

		s.clients[j.Id] = client
	}

	return nil
}
