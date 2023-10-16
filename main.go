package main

import (
	"flag"
	"fmt"
	"log"

	grpc_client "github.com/mrtdeh/centor/pkg/grpc/client"
	grpc_server "github.com/mrtdeh/centor/pkg/grpc/server"
	"github.com/mrtdeh/centor/proto"
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

	name := flag.String("n", "", "server host")
	host := flag.String("h", "0.0.0.0", "server host")
	port := flag.String("p", "10000", "server port")
	isServer := flag.Bool("server", false, "is server")
	serverAddr := flag.String("server-addr", "", "server address for dialing")
	flag.Parse()

	if name == nil {
		log.Fatal("do not set your name")
	}

	if *isServer { // if mode is server

		addr := fmt.Sprintf("%s:%s", *host, *port)
		s := grpc_server.NewServer(*name)
		// listen address for any connection
		err := s.Serve(addr)
		if err != nil {
			log.Fatal(err)
		}

	} else { // if mode is client

		c := grpc_client.NewClient(*name)
		// dialing to server
		if err := c.Dial(*serverAddr); err != nil {
			log.Fatal(err)
		}
		// follow grpc server connction
		if err := c.Follow(); err != nil {
			log.Fatal(err)
		}

	}
}
