package grpc_server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mrtdeh/centor/proto"
)

const (
	ChangeActionAdd = iota
	ChangeActionRemove
)

func (a *agent) applyChange(id string, ni NodeInfo, action int32) error {
	if id == "" {
		return fmt.Errorf("id is empty, must be exist")
	}
	switch action {
	case ChangeActionAdd:
		nodesInfo[id] = ni
	case ChangeActionRemove:
		delete(nodesInfo, id)
	}

	data, err := json.Marshal(nodesInfo)
	if err != nil {
		return err
	}
	for _, child := range a.childs {
		// ignore childs that are leader
		if child.isLeader {
			continue
		}
		// notice to childs
		child.proto.Notice(context.Background(), &proto.NoticeRequest{
			Notice: &proto.NoticeRequest_NodesChange{
				NodesChange: &proto.NodesChange{
					Id:   a.id,
					Data: string(data),
				},
			},
		})
	}

	return nil
}
