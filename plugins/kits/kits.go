package PluginKits

import (
	"context"
	"fmt"

	grpc_server "github.com/mrtdeh/centor/pkg/grpc/server"
)

type CentorHandler interface {
	Call() (string, error)
	GetClusterNodes() map[string]grpc_server.NodeInfo
	WaitForReady(context.Context) error
}

// Plugin interface
type Plugin interface {
	Run()
}

// CoreService structure
type PluginManagerService struct {
	Plugins []Plugin
}

// AddPlugin method to add a plugin to the CoreService
func (c *PluginManagerService) AddPlugin(p Plugin) {
	c.Plugins = append(c.Plugins, p)
}

// Start method to start the CoreService and its plugins
func (c *PluginManagerService) Start() {
	fmt.Println("PluginManagerService is starting...")
	for _, plugin := range c.Plugins {
		go plugin.Run()
	}
}
