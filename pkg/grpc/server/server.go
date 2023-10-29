package grpc_server

import (
	"fmt"
)

type Config struct {
	Name     string
	Host     string
	Port     uint
	Replica  []string
	IsServer bool
	IsLeader bool
}

func Init(cnf Config) error {
	var MainErr chan error

	a := &agent{
		id:       cnf.Name,
		addr:     fmt.Sprintf("%s:%d", cnf.Host, cnf.Port),
		childs:   make(map[string]child),
		done:     make(chan bool, 1),
		isServer: cnf.IsServer,
		isLeader: cnf.IsLeader,
		brothers: cnf.Replica,
	}

	if cnf.IsServer {
		go func() {
			MainErr <- a.Listen()
		}()
	}

	if !cnf.IsLeader && len(cnf.Replica) > 0 {
		go func() {
			MainErr <- a.Connect()
		}()
	}

	return <-MainErr
}
