package api_v1

import (
	"github.com/gin-gonic/gin"
	grpc_server "github.com/mrtdeh/centor/pkg/grpc/server"
)

func GetNodes(c *gin.Context) {
	res := grpc_server.GetClusterNodes()
	var r []any
	for _, v := range res {
		r = append(r, v)
	}
	c.JSON(200, gin.H{
		"result": r,
	})
}
