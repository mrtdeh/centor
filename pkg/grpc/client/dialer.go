package grpc_client

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"github.com/mrtdeh/centor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Configs struct {
	Name                string
	IsServer            bool
	ConnectToOnlyMaster bool
	ServersAddr         []string
}
type client struct {
	id string
	// addr       string
	serverAddr string
	isServer   bool
	conn       proto.DiscoveryClient
}

// connect to server
func (cnf *Configs) Connect() error {
	// master election for servers / best election for clients
	var addr string
	var err error
	if cnf.IsServer {
		addr, err = leaderElect(cnf.ServersAddr)
	} else {
		addr, err = bestElect(cnf.ServersAddr)
	}
	if err != nil {
		return err
	}

	// dial to selected master
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("error in dial : %s", err.Error())
	}
	// create client object
	c := &client{
		id:         cnf.Name,
		serverAddr: addr,
		isServer:   cnf.IsServer,
		conn:       proto.NewDiscoveryClient(conn),
	}

	go c.startServiceProxy()

	// client do follow master
	if err := c.follow(); err != nil {
		return fmt.Errorf("error in follow : %s", err.Error())
	}

	return nil
}

func (c *client) startServiceProxy() {

	c.wait()
	// create proxy request
	if err := c.createProxy("client2", "8090"); err != nil {
		log.Fatalf("error in proxy : %s", err.Error())
	}

	// connect to proxy server (gRPC)
	cc, err := grpc.Dial("localhost:11001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error in dial : %s", err.Error())
	}
	ccc := proto.NewProxyManagerClient(cc)

	listener, err := net.Listen("tcp", ":8001")
	if err != nil {
		panic("connection error:" + err.Error())
	}

	for {
		lc, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept Error:", err)
			continue
		}
		// copyConn(conn)
		connData := make([]byte, 1024)
		_, err = lc.Read(connData)
		if err != nil {
			if err == io.EOF {
				log.Println("read EOF")
				continue
			}
			log.Fatal("error in read connection222 : ", err.Error())
		}
		res, err := ccc.SendPayload(context.Background(), &proto.RequestPayload{
			Conn: connData,
		})
		if err != nil {
			log.Fatal("error in SendPayload : ", err.Error())
		}

		lc.Write(res.Body)
	}

}

func bestElect(addrs []string) (string, error) {

	index := -1
	weight := math.MaxInt32
	for i, a := range addrs {

		conn, err := grpc.Dial(a, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return "", fmt.Errorf("error in dial : %s", err.Error())
		}

		c := proto.NewDiscoveryClient(conn)
		res, err := c.GetInfo(context.Background(), &proto.EmptyRequest{})
		if err != nil {
			// log.Println("error in getInfo : ", err.Error())
			fmt.Println("failed to get info from  :", a)
			continue
		}

		if res.Weight < int32(weight) {
			weight = int(res.Weight)
			index = i
			conn.Close()
		}

	}

	if index == -1 {
		return "", fmt.Errorf("server's are not available")
	}

	return addrs[index], nil

}

func leaderElect(addrs []string) (string, error) {
	for _, a := range addrs {
		conn, err := grpc.Dial(a, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return "", fmt.Errorf("error in dial : %s", err.Error())
		}
		c := proto.NewDiscoveryClient(conn)
		res, err := c.GetInfo(context.Background(), &proto.EmptyRequest{})
		if err != nil {
			return "", fmt.Errorf("error in getInfo : %s", err.Error())
		}

		if res.IsMaster {
			conn.Close()
			return a, nil
		}

	}
	return "", fmt.Errorf("leader not found")
}
