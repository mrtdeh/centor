package grpc_server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mrtdeh/centor/proto"
)

func (a *agent) Notice(ctx context.Context, req *proto.NoticeRequest) (*proto.Close, error) {
	c := &proto.Close{}

	if nch := req.GetNodesChange(); nch != nil {
		fmt.Println("New notice - change nodes : ", nch.Id)
		var nodes map[string]NodeInfo
		err := json.Unmarshal([]byte(nch.Data), &nodes)
		if err != nil {
			return c, err
		}
		nodesInfo = nodes

		if a.childs != nil {
			for _, child := range a.childs {
				child.proto.Notice(context.Background(), &proto.NoticeRequest{
					Notice: &proto.NoticeRequest_NodesChange{
						NodesChange: &proto.NodesChange{
							Data: nch.Data,
						},
					},
				})
			}
		}
	}

	return &proto.Close{}, nil
}
