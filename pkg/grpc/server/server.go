package grpc_server

import (
	"fmt"
	"time"
)

type Config struct {
	Name     string   // Name of the server(id)
	Host     string   // Host of the server
	AltHost  string   // alternative host of the server (optional)
	Port     uint     // Port of the server
	Replica  []string // servers addresses for replication
	IsServer bool     // is this node a server or not
	IsLeader bool     // is this node leader or not
}

var a *agent

func Start(cnf Config) error {
	if cnf.Host == "" {
		cnf.Host = "127.0.0.1"
	}

	// resolve alternative host from config
	var host string = cnf.Host
	if cnf.AltHost != "" {
		host = cnf.AltHost
	}

	// create default agent instance
	a = &agent{
		id:       cnf.Name,
		addr:     fmt.Sprintf("%s:%d", host, cnf.Port),
		childs:   make(map[string]*child),
		isServer: cnf.IsServer,
		isLeader: cnf.IsLeader,
		servers:  cnf.Replica,
	}

	// connect to leader if not a leader and there are servers in the cluster
	if !cnf.IsLeader && len(cnf.Replica) > 0 {
		var err error
		go func() {
			for {
				// try connect to parent server
				err = a.ConnectToParent()
				if err != nil {
					fmt.Println(err.Error())
				}
				// retry delay time 1 second
				time.Sleep(time.Second * 1)
			}
		}()
	} else {
		// if is a leader or there are no servers in the cluster
		// add current node info to nodes info map
		nodesInfo[a.id] = NodeInfo{
			Id:       a.id,
			Address:  a.addr,
			IsServer: true,
			IsLeader: true,
		}
	}

	return a.Listen()
}
