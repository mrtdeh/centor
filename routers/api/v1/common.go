package api_v1

import (
	grpc_server "github.com/mrtdeh/centor/pkg/grpc/server"
)

var h = &grpc_server.CoreHandlers{
	Agent: grpc_server.App,
}

// func getServerAPI() *grpc_server.CoreHandlers {
// 	for {
// 		a := grpc_server.App
// 		if a != nil && a.is{
// 			return &grpc_server.CoreHandlers{
// 				agent: grpc_server.App,
// 			}
// 		}
// 		time.Sleep(time.Millisecond * 500)
// 		log.Println("API Server is not running, waiting...")
// 	}
// }
