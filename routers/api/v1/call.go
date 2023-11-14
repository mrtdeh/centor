package api_v1

import (
	"github.com/gin-gonic/gin"
	grpc_server "github.com/mrtdeh/centor/pkg/grpc/server"
)

func Call(c *gin.Context) {
	tags, err := grpc_server.Call()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"result": tags,
	})
}
