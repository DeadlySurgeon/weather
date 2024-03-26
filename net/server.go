package net

import (
	"fmt"
	"net/http"
	"os"

	"github.com/deadlysurgeon/weather/settings"
	"github.com/deadlysurgeon/weather/weather"
)

// Server is our network server structure. Allows for easy dependency sharing
// and configuration.
type Server struct {
	config  settings.App
	weather weather.Service
	server  *http.Server
}

// New creates a new net server instance.
func New(config settings.App, weather weather.Service) (*Server, error) {
	handler := http.NewServeMux()
	s := &Server{
		config:  config,
		weather: weather,
		server: &http.Server{
			Addr:    config.Net.Bind,
			Handler: handler,
		},
	}

	handler.HandleFunc("/weather/at", s.weatherAt())

	return s, nil
}

// Start starts the server synchronously.
func (s *Server) Start() error {
	if s.config.Net.TLSCert != "" && s.config.Net.TLSKey != "" {
		return s.server.ListenAndServeTLS(s.config.Net.TLSCert, s.config.Net.TLSKey)
	}

	return s.server.ListenAndServe()
}

// Stop stops the server.
func (s *Server) Stop() {
	if err := s.server.Close(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to stop server:", err)
	}
}
