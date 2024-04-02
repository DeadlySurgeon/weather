package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/deadlysurgeon/weather/net"
	"github.com/deadlysurgeon/weather/settings"
	"github.com/deadlysurgeon/weather/weather"
)

const maxAttempts = 10

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	config, err := settings.Process[settings.App]()
	if err != nil {
		return fmt.Errorf("failed to set up config: %w", err)
	}

	weather, err := weather.New(config)
	if err != nil {
		return fmt.Errorf("failed to set up weather service: %w", err)
	}

	server, err := net.New(config, weather)
	if err != nil {
		return fmt.Errorf("failed to set up net server instance: %w", err)
	}

	go notify(server)

	return serverLoop(server)
}

func notify(server *net.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("Stopping server")

	server.Stop()
	<-c
	os.Exit(1)
}

func serverLoop(server *net.Server) error {
	var attempts int
	for {
		if err := server.Start(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			fmt.Fprintln(os.Stderr, "Server returned an error:", err)
		}
		attempts++
		if attempts >= maxAttempts {
			return fmt.Errorf("Failed too many times")
		}
		fmt.Println("Attempting to restart server...")
	}
}
