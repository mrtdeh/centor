package grpc_server

import (
	"context"
	"strings"

	"github.com/mrtdeh/centor/proto"
)

func Call() (string, error) {

	res, err := a.parent.proto.Call(context.Background(), &proto.CallRequest{
		AgentId: a.id,
	})
	if err != nil {
		return "", err
	}
	return strings.Join(res.Tags, " ,"), nil
}

func GetClusterNodes() map[string]NodeInfo {
	return nodesInfo
}
