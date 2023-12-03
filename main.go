package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/common-nighthawk/go-figure"
	api_server "github.com/mrtdeh/centor/pkg/api"
	"github.com/mrtdeh/centor/pkg/config"
	grpc_server "github.com/mrtdeh/centor/pkg/grpc/server"
	pluginManager "github.com/mrtdeh/centor/plugins"
	PluginKits "github.com/mrtdeh/centor/plugins/assets"
	"github.com/mrtdeh/centor/routers"
)

func main() {

	// print centor in cli
	printLogo()

	// load configurations
	cnf := config.LoadConfiguration()

	var serversAddrs []string
	sd := cnf.ServersAddr
	if sd != "" {
		serversAddrs = strings.Split(strings.TrimSpace(sd), ",")
	}

	var primariesAddrs []string
	pd := cnf.PrimaryServersAddr
	if pd != "" {
		primariesAddrs = strings.Split(strings.TrimSpace(pd), ",")
	}

	// initilize api server
	router := routers.InitRouter()

	// bootstrap plugins
	err := pluginManager.Bootstrap(pluginManager.Config{
		Config: PluginKits.Config{
			GRPCHandler: &grpc_server.GRPC_Handlers{},
			RouterAPI:   router,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// start api server
	if config.WithAPI {
		httpServer := api_server.HttpServer{
			Host:   "0.0.0.0",
			Port:   9090,
			Router: router,
		}
		fmt.Printf("initil api server an address %s:%d\n", httpServer.Host, httpServer.Port)

		go func() {
			log.Fatal(httpServer.Serve())
		}()
	}

	// start gRPC server
	err = grpc_server.Start(grpc_server.Config{
		Name:       cnf.Name,
		DataCenter: cnf.DataCenter,
		Host:       cnf.Host,
		AltHost:    cnf.AltHost,
		Port:       cnf.Port,
		IsServer:   cnf.IsServer,
		IsLeader:   cnf.IsLeader,
		Servers:    serversAddrs,
		Primaries:  primariesAddrs,
	})
	if err != nil {
		log.Fatal(err)
	}

}

func printLogo() {
	myFigure := figure.NewFigure("CENTOR", "", true)
	myFigure.Print()
	fmt.Println()
}
