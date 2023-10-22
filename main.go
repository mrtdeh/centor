package main

import (
	"flag"
	"log"
	"strings"

	"github.com/mrtdeh/centor/pkg/config"
	client "github.com/mrtdeh/centor/pkg/grpc/client"
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

	if cnf.IsServer { // if mode is server

		s := server.Configs{
			Name:     cnf.Name,
			Host:     cnf.Host,
			Port:     cnf.Port,
			IsLeader: cnf.IsLeader,
			Replica:  serversAddrs,
		}
		err := s.Listen()
		if err != nil {
			log.Fatal(err)
		}

	} else { // if mode is client

		c := client.Configs{
			Name:        cnf.Name,
			ServersAddr: serversAddrs,
		}
		if err := c.Connect(); err != nil {
			log.Fatal(err)
		}

	}
}
