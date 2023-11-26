package echo_plugin

import (
	"context"
	"fmt"

	PluginKits "github.com/mrtdeh/centor/plugins/kits"
)

// ExamplePlugin1 is an example plugin implementing the Plugin interface
type PluginProvider struct {
	Name    string
	Handler PluginKits.CentorHandler
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

	tags, err := h.Call()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("from echo plugin: ", tags)
}
