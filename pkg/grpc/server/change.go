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

type NodeInfo struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Port     string `json:"port"`
	IsServer bool   `json:"is_server"`
	IsLeader bool   `json:"is_leader"`
	ParentId string `json:"parent_id"`
}

const (
	ChangeActionAdd = iota
	ChangeActionRemove
)

func (a *agent) Change(ctx context.Context, req *proto.ChangeRequest) (*proto.Close, error) {
	c := &proto.Close{}
	if !a.isLeader {
		return c, fmt.Errorf("you must send change request to leader not here")
	}

	if nch := req.GetNodesChange(); nch != nil {
		// fmt.Println("New change - change nodes")

		var ni NodeInfo
		err := json.Unmarshal([]byte(nch.Data), &ni)
		if err != nil {
			return c, err
		}

		err = a.applyChange(nch.Id, ni, nch.Action)
		if err != nil {
			return c, err
		}
	}

	return c, nil
}

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
