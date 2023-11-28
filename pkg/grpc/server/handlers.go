package grpc_server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/mrtdeh/centor/proto"
)

type GRPC_Handlers struct{}

type FileHandler struct {
	Name      string
	Extension string
	Data      []byte
}

func (h *GRPC_Handlers) SendFile(ctx context.Context, nodeId, filename string, data []byte) error {

	reader := bytes.NewReader(data)
	filesize := reader.Size()
	buffer := make([]byte, 1024)

	if n, ok := nodesInfo[nodeId]; ok {
		// fmt.Printf("send to : %+v\n", n)
		conn, err := grpc_Dial(n.Address)
		if err != nil {
			return err
		}
		defer func() {
			time.Sleep(time.Second)
			conn.Close()
			fmt.Println("closed connection")
		}()

		client := proto.NewDiscoveryClient(conn)
		stream, err := client.SendFile(context.Background())
		if err != nil {
			return err
		}

		firstInfo := &proto.SendFileRequest{
			Data: &proto.SendFileRequest_Info{
				Info: &proto.FileInfo{
					Name: filename,
					Size: uint32(filesize),
				},
			},
		}
		err = stream.Send(firstInfo)
		if err != nil {
			fmt.Println("debug send error : ", err)
			return err
		}

		fmt.Println("debug send 1")

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
			fmt.Println("debug send 2")
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
		fmt.Println("debug send 3")

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
