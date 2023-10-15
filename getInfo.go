package main

import (
	"context"

	"github.com/mrtdeh/centor/proto"
)

func (s *Server) GetInfo(context.Context, *proto.EmptyRequest) (*proto.InfoResponse, error) {
	return &proto.InfoResponse{
		Id:       s.Id,
		IsMaster: s.IsMaster,
	}, nil
}
