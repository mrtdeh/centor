package grpc_server

import (
	"fmt"
	"time"
)

type Config struct {
	Name       string   // Name of the server(id)
	DataCenter string   // Name of the server(id)
	Host       string   // Host of the server
	AltHost    string   // alternative host of the server (optional)
	Port       uint     // Port of the server
	Servers    []string // servers addresses for replication
	Primaries  []string // primaries addresses
	IsServer   bool     // is this node a server or not
	IsLeader   bool     // is this node leader or not
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
		dc:       cnf.DataCenter,
		addr:     fmt.Sprintf("%s:%d", host, cnf.Port),
		childs:   make(map[string]*child),
		isServer: cnf.IsServer,
		isLeader: cnf.IsLeader,
	}

	var finishCh = make(chan struct{}, 1)
	go func() {
		select {
		case <-finishCh:
			a.isReady = true
			return
		}
	}()

	// *** CONNECT TO LEADER SERVER *** if this node is not a leader
	if !cnf.IsLeader && len(cnf.Servers) > 0 {
		var err error
		go func() {
			for {
				// try connect to leader server
				err = a.ConnectToParent(connectConfig{
					ServersAddresses: cnf.Servers,
					OnFinishChan:     finishCh,
				})
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
			Id:         a.id,
			Address:    a.addr,
			IsServer:   true,
			IsLeader:   true,
			DataCenter: a.dc,
		}

		// *** CONNECT TO PRIMARY SERVER *** if the primary server is available
		if len(cnf.Primaries) > 0 {
			var err error
			go func() {
				for {
					// try connect to primary server
					err = a.ConnectToParent(connectConfig{
						ConnectToPrimary: true,
						ServersAddresses: cnf.Primaries,
						OnFinishChan:     finishCh,
					})
					if err != nil {
						fmt.Println(err.Error())
					}
					// retry delay time 1 second
					time.Sleep(time.Second * 1)
				}
			}()
		}

	}

	fmt.Println("DataCenter : ", cnf.DataCenter)
	return a.Listen()
}
