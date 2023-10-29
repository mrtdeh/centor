package main

import (
	"flag"
	"log"
	"strings"

	"github.com/mrtdeh/centor/pkg/config"
	server "github.com/mrtdeh/centor/pkg/grpc/server"
)

func main() {
	var confPath string
	flag.StringVar(&confPath, "c", "", "config path")
	flag.Parse()

	cnf := config.LoadConfiguration(confPath)

	var serversAddrs []string = nil
	sd := cnf.ServersAddr
	if sd != "" {
		serversAddrs = strings.Split(strings.TrimSpace(sd), ",")
	}
	// need to implment Start() to run for both server/client
	err := server.Init(server.Config{
		Name:     cnf.Name,
		Host:     cnf.Host,
		Port:     cnf.Port,
		IsServer: cnf.IsServer,
		IsLeader: cnf.IsLeader,
		Replica:  serversAddrs,
	})
	if err != nil {
		log.Fatal(err)
	}

}
