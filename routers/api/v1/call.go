package api_v1

import (
	"github.com/gin-gonic/gin"
)

func Call(c *gin.Context) {
	tags, err := h.Call()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"result": tags,
	})
}
