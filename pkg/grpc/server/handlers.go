package grpc_server

import (
	"context"
	"strings"
	"time"

	"github.com/mrtdeh/centor/proto"
)

type GRPC_Handlers struct{}

func (h *GRPC_Handlers) WaitForReady(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if a.isReady {
				return nil
			}
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func (h *GRPC_Handlers) Call() (string, error) {

	res, err := a.Call(context.Background(), &proto.CallRequest{
		AgentId: a.id,
	})
	if err != nil {
		return "", err
	}
	return strings.Join(res.Tags, " ,"), nil
}

// GetClusterNodes returns a map of all the nodes in the cluster
func (h *GRPC_Handlers) GetClusterNodes() map[string]NodeInfo {
	return nodesInfo
}
