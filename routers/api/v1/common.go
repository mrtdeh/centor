package api_v1

import (
	"log"
	"time"

	grpc_server "github.com/mrtdeh/centor/pkg/grpc/server"
)

func getServerAPI() *grpc_server.CoreHandlers {
	for {
		a := grpc_server.GetAgentInstance()
		if a != nil {
			return &grpc_server.CoreHandlers{
				Agent: grpc_server.GetAgentInstance(),
			}
		}
		time.Sleep(time.Millisecond * 500)
		log.Println("API Server is not running, waiting...")
	}
}
