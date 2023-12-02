package echo_plugin

import (
	"context"
	"fmt"
	"net/http"

	PluginKits "github.com/mrtdeh/centor/plugins/assets"
)

// ExamplePlugin1 is an example plugin implementing the Plugin interface
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

	return nil

}

// Run method for ExamplePlugin1
func (p *PluginProvider) Run() {
	fmt.Printf("Plugin %s is running...\n", p.Name)
	h := p.Handler

	err := h.WaitForReady(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	tags, err := h.Call(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("from echo plugin: ", tags)
}
