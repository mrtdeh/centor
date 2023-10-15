package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	Id       string
	IsMaster bool
	clients  map[string]Client
}
type Client struct {
	Id   string
	Addr string
	conn proto.Discovery_FollowServer
}

var server *Server
var portNumber int = 9001

func main() {

	isServer := flag.Bool("server", false, "is server")
	serverAddr := flag.String("server-addr", "", "is server")
	flag.Parse()

	if *isServer {
		addr := "localhost:90001"
		grpcServer := grpc.NewServer()

		listener, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatalf("error creating the server %v", err)
		}

		proto.RegisterDiscoveryServer(grpcServer, server)
		grpcServer.Serve(listener)
	} else {
		addr := *serverAddr
		id := sha256.Sum256([]byte(time.Now().String()))
		serverId := hex.EncodeToString(id[:])[:6]

		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal(fmt.Errorf("error in dial : %s", err.Error()))
		}
		c := proto.NewDiscoveryClient(conn)

		// im, err := c.GetInfo(context.Background(), &proto.EmptyRequest{})
		// if err != nil {
		// 	log.Fatal(fmt.Errorf("error check is_master : %s\n", err.Error()))
		// }

		// if !im.IsMaster {
		// 	log.Fatal("this server is not master")
		// }

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

				// portNumber++
				// serverPort := fmt.Sprintf("%d", portNumber)

				// service := ServiceProxy{
				// 	// user values
				// 	Name:             "test-service",
				// 	LocalServicePort: p.ProxyPort,
				// 	// protected values
				// 	gRPCServerPort: serverPort,
				// 	msgCh:          make(chan []byte, 1024),
				// }

				// // serviceMap[p.]

				// go func() {
				// 	service.RunProxy()
				// }()
			}

		}
	}
}
