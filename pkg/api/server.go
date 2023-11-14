package api_server

import (
	"fmt"
	"net/http"
	"time"
)

type HttpServer struct {
	Host   string
	Port   uint
	Debug  bool
	Router http.Handler
}

func (s *HttpServer) Serve() error {

	endPoint := fmt.Sprintf("0.0.0.0:%d", s.Port) // Debug Mode
	if !s.Debug {                                 // Release Mode
		endPoint = fmt.Sprintf("%s:%d", s.Host, s.Port)
	}

	server := &http.Server{
		Addr:        endPoint,
		Handler:     s.Router,
		IdleTimeout: time.Second * 10,
	}

	if err := server.ListenAndServe(); err != nil {
		return fmt.Errorf("error in server listener : %w", err)
	}

	return nil
}
