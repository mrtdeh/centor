package pluginManager

import (
	"fmt"

	echo_plugin "github.com/mrtdeh/centor/plugins/echo"
	PluginKits "github.com/mrtdeh/centor/plugins/kits"
)

type Config struct {
	GRPCHandler PluginKits.CentorHandler
}

func Load(cnf Config) error {
	if cnf.GRPCHandler == nil {
		return fmt.Errorf("handler is nil")
	}
	// Create a new instance of pman
	pman := &PluginKits.PluginManagerService{}

	// Add plugins to the pman
	pman.AddPlugin(&echo_plugin.PluginProvider{Name: "echo", Handler: cnf.GRPCHandler})

	// Start the pman and its plugins
	pman.Start()
	return nil
}
