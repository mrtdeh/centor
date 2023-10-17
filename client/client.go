package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var portNumber int = 9001
var serviceProxyMap = map[string]ServiceProxy{}
var serviceMap = map[string]Service{}

type Config struct {
	Service *ServiceConfig `hcl:"service,block"`
}

type ServiceConfig struct {
	Id   string `hcl:"id"`
	Name string `hcl:"name"`
	Port uint   `hcl:"port"`
}

func main() {
	addr := "localhost:90001"

	configs, err := loadConfigsFromDir("/etc/centor.d/")
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}
	parseServices(configs)

	fmt.Printf("services : %+v\n", serviceMap)

	// =============================================================

	id := sha256.Sum256([]byte(time.Now().String()))
	serverId := hex.EncodeToString(id[:])[:6]

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(fmt.Errorf("error in dial : %s", err.Error()))
	}
	c := proto.NewDiscoveryClient(conn)

	req, err := c.Follow(context.Background())
	if err != nil {
		log.Fatalf("failed to follow : %s\n", err.Error())
	}

	req.Send(&proto.FollowerRequest{
		Data: &proto.FollowerRequest_JoinMsg{
			JoinMsg: &proto.JoinMessage{
				Id:   serverId,
				Addr: addr,
			},
		},
	})

	for {
		res, err := req.Recv()
		if err != nil {
			log.Fatalf("error in recivce : %s", err.Error())
		}

		if p := res.GetProxyRequest(); p != nil {

			portNumber++
			serverPort := fmt.Sprintf("%d", portNumber)

			service := ServiceProxy{
				// user values
				Name:             "test-service",
				LocalServicePort: p.ProxyPort,
				// protected values
				gRPCServerPort: serverPort,
				msgCh:          make(chan []byte, 1024),
			}

			// serviceMap[p.]

			go func() {
				service.RunProxy()
			}()
		}

	}
}
