package exec_plugin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	PluginKits "github.com/mrtdeh/centor/plugins/assets"
)

type PluginProvider struct {
	PluginKits.PluginProps
}

func (p *PluginProvider) SetHandler(h PluginKits.CentorHandler) {
	p.Handler = h
}

func (p *PluginProvider) SetRouter(r http.Handler) {
	p.Router = r
}

func (p *PluginProvider) Init() error {
	r, ok := p.Router.(*gin.Engine)
	if !ok {
		return fmt.Errorf("router is not a gin router")
	}
	r.POST("/exec", exec)

	p.Router = r
	return nil

}

var h PluginKits.CentorHandler

// Run method for ExamplePlugin1
func (p *PluginProvider) Run() {
	fmt.Printf("Plugin %s is running...\n", p.Name)
	h = p.Handler

	err := h.WaitForReady(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

}

type ExecRequest struct {
	Command string `json:"command"`
	NodeId  string `json:"node_id"`
}

func exec(c *gin.Context) {
	var req ExecRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Exec(context.Background(), req.NodeId, req.Command)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Exec request on %s : %s\n", req.NodeId, res)
	c.JSON(200, "ok")
}
