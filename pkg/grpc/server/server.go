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

func Init(cnf Config) error {
	var MainErr chan error = make(chan error, 1)

	a := &agent{
		id:       cnf.Name,
		addr:     fmt.Sprintf("%s:%d", cnf.Host, cnf.Port),
		childs:   make(map[string]child),
		isServer: cnf.IsServer,
		isLeader: cnf.IsLeader,
		brothers: cnf.Replica,
		isReady:  false,
	}

	if cnf.IsServer {
		go func() {
			MainErr <- a.Listen()
		}()
	}

	if !cnf.IsLeader && len(cnf.Replica) > 0 {
		var err error
		var try int = 10
		go func() {
			for {
				err = a.Connect()
				if try <= 0 {
					MainErr <- err
				}

				try--
				time.Sleep(time.Second * 1)
			}
		}()
	}

	return <-MainErr
}
