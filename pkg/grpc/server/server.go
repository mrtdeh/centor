package grpc_server

import (
	"fmt"
	"time"
)

type Config struct {
	Name     string
	Host     string
	Port     uint
	Replica  []string
	IsServer bool
	IsLeader bool
}

var a *agent

func Start(cnf Config) error {
	if cnf.Host == "" {
		cnf.Host = "127.0.0.1"
	}

	a = &agent{
		id:       cnf.Name,
		addr:     fmt.Sprintf("%s:%d", cnf.Host, cnf.Port),
		childs:   make(map[string]*child),
		isServer: cnf.IsServer,
		isLeader: cnf.IsLeader,
		servers:  cnf.Replica,
	}

	if !cnf.IsLeader && len(cnf.Replica) > 0 {
		var err error
		go func() {
			for {
				// try connect to parent server
				err = a.ConnectToParent()
				if err != nil {
					fmt.Println(err.Error())
				}
				// retry delay time
				time.Sleep(time.Second * 1)
			}
		}()
	} else {
		nodesInfo[a.id] = NodeInfo{
			Id:       a.id,
			Address:  a.addr,
			IsServer: true,
			IsLeader: true,
		}
	}

	return a.Listen()
}
