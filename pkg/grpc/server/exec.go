package grpc_server

import (
	"context"
	"os/exec"
	"strings"

	"github.com/mrtdeh/centor/proto"
)

func (s *agent) Exec(ctx context.Context, req *proto.ExecRequest) (*proto.ExecResponse, error) {
	cmds := strings.Split(req.Command, " ")
	cmd := exec.Command(cmds[0], cmds[1:]...)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return &proto.ExecResponse{Output: string(out)}, nil
}
