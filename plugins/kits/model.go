package PluginKits

import (
	"net/http"
)

type Config struct {
	GRPCHandler CentorHandler
	RouterAPI   http.Handler
}
