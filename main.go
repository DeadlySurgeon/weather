package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/deadlysurgeon/weather/net"
	"github.com/deadlysurgeon/weather/settings"
	"github.com/deadlysurgeon/weather/weather"
)

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

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		fmt.Println("Stoping server")

		server.Stop()
	}()

	return server.Start()
}
