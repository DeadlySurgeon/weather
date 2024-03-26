package settings

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// App Global Configuration
type App struct {
	Weather Weather
	Net     Net
}

// Net based configuration
type Net struct {
	Bind    string
	TLSCert string
	TLSKey  string
}

// Weather Service based configuration
type Weather struct {
	APIKey   string
	Endpoint string
}

// Process loads the settings.
func Process[T any]() (T, error) {
	var app T

	if err := godotenv.Load(); err != nil {
		// We don't care if an .env is missing, it will be missing in prod.
		if !os.IsNotExist(err) {
			return app, err
		}
	}

	if err := envconfig.Process("", &app); err != nil {
		return app, err
	}

	return app, nil
}
