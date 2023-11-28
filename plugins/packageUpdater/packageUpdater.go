package packageupdater_plugin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	PluginKits "github.com/mrtdeh/centor/plugins/kits"
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
	r.POST("/send-file", sendFile)

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

type SendFileRequest struct {
	Filename string `json:"filename"`
	Data     string `json:"data"`
	NodeId   string `json:"node_id"`
}

func sendFile(c *gin.Context) {
	var req SendFileRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.SendFile(context.Background(), req.NodeId, req.Filename, []byte(req.Data))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, "ok")
}
