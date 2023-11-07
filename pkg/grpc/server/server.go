package grpc_server

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mrtdeh/centor/proto"
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
	if cnf.Host == "" {
		cnf.Host = "127.0.0.1"
	}

	a := &agent{
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
				err = a.ConnectToParent()
				if err != nil {
					fmt.Println(err.Error())
				}
				time.Sleep(time.Second * 1)
			}
		}()
	}

	if cnf.Name == "test-call" {
		go func() {
			time.Sleep(time.Second * 3)
			res, err := a.parent.proto.Call(context.Background(), &proto.CallRequest{
				AgentId: a.id,
			})
			if err != nil {
				log.Fatal("error in call : ", err.Error())
			}
			log.Println("tags : ", res.Tags)
		}()
	}

	return a.Listen()
}
