package main

import (
	"log"
	"strings"

	api_server "github.com/mrtdeh/centor/pkg/api"
	"github.com/mrtdeh/centor/pkg/config"
	grpc_server "github.com/mrtdeh/centor/pkg/grpc/server"
	"github.com/mrtdeh/centor/routers"
)

func main() {

	cnf := config.LoadConfiguration()

	var serversAddrs []string
	sd := cnf.ServersAddr
	if sd != "" {
		serversAddrs = strings.Split(strings.TrimSpace(sd), ",")
	}

	if config.WithAPI {
		httpServer := api_server.HttpServer{
			Host:   "localhost",
			Port:   9090,
			Debug:  true,
			Router: routers.InitRouter(),
		}

		go func() {
			log.Fatal(httpServer.Serve())
		}()
	}

	err := grpc_server.Start(grpc_server.Config{
		Name:     cnf.Name,
		Host:     cnf.Host,
		AltHost:  cnf.AltHost,
		Port:     cnf.Port,
		IsServer: cnf.IsServer,
		IsLeader: cnf.IsLeader,
		Replica:  serversAddrs,
	})
	if err != nil {
		log.Fatal(err)
	}

}
