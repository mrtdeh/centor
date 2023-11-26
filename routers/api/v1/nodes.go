package api_v1

import (
	"github.com/gin-gonic/gin"
)

func GetNodes(c *gin.Context) {
	res := h.GetClusterNodes()
	var r []any
	for _, v := range res {
		r = append(r, v)
	}
	c.JSON(200, gin.H{
		"result": r,
	})
}
