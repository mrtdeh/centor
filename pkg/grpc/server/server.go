package grpc_server

import (
	"fmt"
	"log"
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

func Start(cnf Config) error {
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

	if !cnf.IsLeader && len(cnf.Replica) > 0 {
		var err error
		go func() {
			for {
				err = a.Connect()
				if err != nil {
					log.Println("failed to connect to server : ", err.Error())
				}
				time.Sleep(time.Second * 1)
			}
		}()
	}

	if cnf.IsServer {
		MainErr <- a.Listen()
	}

	return <-MainErr
}
