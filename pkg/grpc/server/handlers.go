package grpc_server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mrtdeh/centor/proto"
)

type GRPC_Handlers struct{}

func (h *GRPC_Handlers) SendFile(ctx context.Context, nodeId, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("cannot open txt file: ", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	stat, _ := file.Stat()

	if n, ok := nodesInfo[nodeId]; ok {
		conn, err := grpc_Dial(n.Address)
		if err != nil {
			return err
		}
		defer conn.Close()

		client := proto.NewDiscoveryClient(conn)
		stream, err := client.SendFile(ctx)
		if err != nil {
			return err
		}
		defer stream.CloseSend()

		firstInfo := &proto.SendFileRequest{
			Data: &proto.SendFileRequest_Info{
				Info: &proto.FileInfo{
					Name: stat.Name(),
					Size: uint32(stat.Size()),
				},
			},
		}
		err = stream.Send(firstInfo)
		if err != nil {
			return err
		}

		for {
			n, err := reader.Read(buffer)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal("cannot read chunk to buffer: ", err)
			}

			chunk := &proto.SendFileRequest{
				Data: &proto.SendFileRequest_ChunkData{
					ChunkData: buffer[:n],
				},
			}
			err = stream.Send(chunk)
			if err != nil {
				return err
			}
		}

		endInfo := &proto.SendFileRequest{
			Data: &proto.SendFileRequest_End{
				End: true,
			},
		}
		err = stream.Send(endInfo)
		if err != nil {
			return err
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
