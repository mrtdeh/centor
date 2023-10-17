package main

import (
	"log"

	"github.com/mrtdeh/centor/pkg/config"
	client "github.com/mrtdeh/centor/pkg/grpc/client"
	server "github.com/mrtdeh/centor/pkg/grpc/server"
)

func main() {

	cnf := config.LoadConfiguration()

	if *cnf.IsServer { // if mode is server

		s := server.Configs{
			Name: *cnf.Name,
			Host: *cnf.Host,
			Port: *cnf.Port,
		}
		err := s.Listen()
		if err != nil {
			log.Fatal(err)
		}

	} else { // if mode is client

		c := client.Configs{
			Name: *cnf.Name,
			ServersAddr: []string{
				*cnf.ServerAddr,
			},
		}
		if err := c.Connect(); err != nil {
			log.Fatal(err)
		}

	}
}
