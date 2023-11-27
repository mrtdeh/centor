package PluginKits

import "fmt"

func validateConfig(config Config) error {
	if config.GRPCHandler == nil {
		return fmt.Errorf("handler is nil")
	}
	if config.RouterAPI == nil {
		return fmt.Errorf("router API is nil")
	}
	return nil
}
