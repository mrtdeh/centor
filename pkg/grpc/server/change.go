package grpc_server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mrtdeh/centor/proto"
)

var (
	nodesInfo = map[string]NodeInfo{}
)

func (a *agent) Change(ctx context.Context, req *proto.ChangeRequest) (*proto.Close, error) {
	c := &proto.Close{}
	if !a.isLeader {
		return c, fmt.Errorf("you must send change request to primary not here")
	}

	if nch := req.GetNodesChange(); nch != nil {
		// fmt.Println("New change - change nodes")

		var ni NodeInfo
		err := json.Unmarshal([]byte(nch.Data), &ni)
		if err != nil {
			return c, err
		}

		err = a.syncAgentChange(&agent{
			id:        ni.Id,
			addr:      ni.Address,
			isServer:  ni.IsServer,
			isLeader:  ni.IsLeader,
			isPrimary: ni.IsPrimary,
			dc:        ni.DataCenter,
			parent:    &parent{agent: agent{id: ni.ParentId}},
		}, nch.Action)
		if err != nil {
			return c, err
		}

	}

	return c, nil
}
