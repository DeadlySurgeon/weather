package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/deadlysurgeon/weather/settings"
)

type impl struct {
	endpoint string
	apiKey   string
	client   *http.Client
}

// New creates a weather service implementation
func New(config settings.App) (Service, error) {
	service := &impl{
		endpoint: "https://api.openweathermap.org/data/3.0/onecall", // ?lat={lat}&lon={lon}&exclude={part}&appid={API key}
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}

	if service.apiKey = config.Weather.APIKey; service.apiKey == "" {
		return nil, fmt.Errorf("no api key specified")
	}

	if config.Weather.Endpoint != "" {
		service.endpoint = config.Weather.Endpoint
	}

	return service, nil
}

func (s *impl) At(lat, lon string) (Report, error) {
	report := Report{
		Temperature:    "unknown",
		Condition:      "unknown",
		TemperatureRaw: 0,
	}

	req, err := s.formRequest(lat, lon)
	if err != nil {
		return report, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return report, err
	}
	defer resp.Body.Close()

	// TODO:
	// - Double check how they return errors such as overuse.

	if resp.StatusCode != http.StatusOK {
		return report, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	var owmResp openweathermapResponse
	if err = json.NewDecoder(resp.Body).Decode(&owmResp); err != nil {
		return report, fmt.Errorf("failed to read response: %w", err)
	}

	report.TemperatureRaw = owmResp.Current.FeelsLike
	report.Temperature = feelsLike(report.TemperatureRaw)
	if len(owmResp.Current.Weather) > 0 {
		report.Condition = owmResp.Current.Weather[0].Main
	}

	return report, nil
}

func feelsLike(temp float32) string {
	switch {
	case temp <= 0:
		return "freezing"
	case temp <= 10:
		return "cold"
	case temp <= 20:
		return "moderate"
	case temp <= 30:
		return "hot"
	default:
		return "burning"
	}
}

func (s *impl) formRequest(lat, lon string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, s.endpoint, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("appid", s.apiKey)
	q.Add("units", "metric")
	q.Add("lat", lat)
	q.Add("lon", lon)
	req.URL.RawQuery = q.Encode()

	return req, nil
}

// Only grab the fields we care about. If we needed more from the service we
// would need to break this out into its own models file and each sub struct be
// exportable.
type openweathermapResponse struct {
	Current struct {
		FeelsLike float32 `json:"feels_like"`
		Weather   []struct {
			Main string `json:"main"`
		} `json:"weather"`
	} `json:"current"`
}
