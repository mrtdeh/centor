package grpc_client

import (
	"context"
	"fmt"

	"github.com/mrtdeh/centor/proto"
)

func (c *client) createProxy(target, port string) error {
	_, err := c.conn.CreateProxy(context.Background(), &proto.CreateProxyRequest{
		TargetId:          target,
		TargetServicePort: port,
	})
	if err != nil {
		return fmt.Errorf("failed to proxy : %s\n", err.Error())
	}

	return nil
}
