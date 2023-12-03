package grpc_server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/mrtdeh/centor/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type GRPC_Handlers struct{}

type FileHandler struct {
	Name      string
	Extension string
	Data      []byte
}

func (h *GRPC_Handlers) GetMyId() string {
	return a.id
}

func (h *GRPC_Handlers) GetParentId() string {
	if a.parent != nil {
		return a.parent.id
	}
	return ""
}

func (h *GRPC_Handlers) FireEvent(ctx context.Context, nodeId, event string, params ...any) error {
	protoParams := []*anypb.Any{}
	for _, p := range params {
		anyValue, err := ConvertInterfaceToAny(p)
		if err != nil {
			return err
		}
		protoParams = append(protoParams, anyValue)
	}
	// check if node_id is exist or not
	if n, ok := nodesInfo[nodeId]; ok {
		conn, err := grpc_Dial(n.Address)
		if err != nil {
			return err
		}
		defer conn.Close()

		client := proto.NewDiscoveryClient(conn)
		_, err = client.FireEvent(context.Background(), &proto.EventRequest{
			Name:   event,
			Params: protoParams,
		})
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (h *GRPC_Handlers) Exec(ctx context.Context, nodeId, commnad string) (string, error) {

	// check if node_id is exist or not
	if n, ok := nodesInfo[nodeId]; ok {
		conn, err := grpc_Dial(n.Address)
		if err != nil {
			return "", err
		}
		defer conn.Close()

		client := proto.NewDiscoveryClient(conn)
		res, err := client.Exec(context.Background(), &proto.ExecRequest{
			Command: commnad,
		})
		if err != nil {
			return "", err
		}

		return res.Output, nil
	}

	return "", nil
}

func (h *GRPC_Handlers) SendFile(ctx context.Context, nodeId, filename string, data []byte) error {

	reader := bytes.NewReader(data)
	filesize := reader.Size()
	buffer := make([]byte, 1024)

	// check if node_id is exist or not
	if n, ok := nodesInfo[nodeId]; ok {
		conn, err := grpc_Dial(n.Address)
		if err != nil {
			return err
		}
		defer conn.Close()

		client := proto.NewDiscoveryClient(conn)
		stream, err := client.SendFile(context.Background())
		if err != nil {
			return err
		}

		// send the file information
		err = stream.Send(&proto.SendFileRequest{
			Data: &proto.SendFileRequest_Info{
				Info: &proto.FileInfo{
					Name: filename,
					Size: uint32(filesize),
				},
			},
		})
		if err != nil {
			return err
		}

		for {
			n, err := reader.Read(buffer)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			// send the chunks of the file
			err = stream.Send(&proto.SendFileRequest{
				Data: &proto.SendFileRequest_ChunkData{
					ChunkData: buffer[:n],
				},
			})
			if err != nil {
				return err
			}
		}

		// send the end of the file
		err = stream.Send(&proto.SendFileRequest{
			Data: &proto.SendFileRequest_End{
				End: true,
			},
		})
		if err != nil {
			return err
		}

		// receive server response and error if any
		_, err = stream.CloseAndRecv()
		if err != nil && err != io.EOF {
			return fmt.Errorf("Error receiving response: %v", err)
		}

	} else {
		return fmt.Errorf("Node %s not found", nodeId)
	}

	return nil
}

// wait for current agent is running completely
func (h *GRPC_Handlers) WaitForReady(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if a != nil && a.isReady {
				return nil
			}
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func (h *GRPC_Handlers) Call(ctx context.Context) (string, error) {

	res, err := a.Call(ctx, &proto.CallRequest{
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
