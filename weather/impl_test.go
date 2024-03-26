package weather

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFeelsLike(t *testing.T) {
	for name, test := range map[string]float32{
		"freezing": -4,
		"cold":     5,
		"moderate": 11,
		"hot":      22,
		"burning":  42,
	} {
		test := test
		t.Run(name, func(t *testing.T) {
			if s := feelsLike(test); s != name {
				t.Fatalf("Expected %s got %s", name, s)
			}
		})
	}
}

func TestFormRequest(t *testing.T) {
	for name, test := range map[string]struct {
		url         string
		expectError bool
	}{
		"good url": {
			url: "https://example.com",
		},
		"bad url": {
			url:         "http://%41:8080/",
			expectError: true,
		},
	} {
		test := test
		t.Run(name, func(t *testing.T) {
			_ = test
			_, err := (&impl{endpoint: test.url}).formRequest("", "")
			if (err != nil) != test.expectError {
				t.Fatalf("Expected error (%v) got %v", test.expectError, err)
			}
		})
	}
}

func TestWeatherAt(t *testing.T) {
	for name, test := range map[string]struct {
		lat, lon       string
		mockResponse   string
		endpoint       string
		mockStatusCode int
		clientErr      error
		wantErr        bool
		wantReport     Report
	}{
		"successful report": {
			lat:            "0",
			lon:            "0",
			mockResponse:   `{"current":{"feels_like":25.3,"weather":[{"main":"Clear"}]}}`,
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			wantReport: Report{
				Condition:      "Clear",
				Temperature:    "hot",
				TemperatureRaw: 25.3,
			},
		},
		"bad status code": {
			lat:            "invalid",
			lon:            "invalid",
			mockResponse:   `Bad Request`,
			mockStatusCode: http.StatusBadRequest,
			wantErr:        true,
		},
		"invalid JSON response": {
			lat:            "0",
			lon:            "0",
			mockResponse:   `Invalid JSON`,
			mockStatusCode: http.StatusOK,
			wantErr:        true,
		},
		"do client error": {
			lat:            "0",
			lon:            "0",
			clientErr:      fmt.Errorf("client error"),
			mockStatusCode: http.StatusOK,
			wantErr:        true,
		},
		"bad endpoint": {
			lat:      "0",
			lon:      "0",
			endpoint: "http://%41:8080/",
			wantErr:  true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			// Setup mock HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.mockStatusCode)
				fmt.Fprintln(w, test.mockResponse)
			}))
			defer server.Close()

			// Create an instance of `impl` with a mocked client
			s := &impl{
				endpoint: test.endpoint,
				client: mockHTTPClient(func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode: test.mockStatusCode,
						Body:       io.NopCloser(strings.NewReader(test.mockResponse)),
						Header:     make(http.Header),
					}
				}, test.clientErr),
			}

			// Execute WeatherAt
			gotReport, err := s.At(test.lat, test.lon)
			if (err != nil) != test.wantErr {
				t.Errorf("WeatherAt() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if !test.wantErr && !compareReports(gotReport, test.wantReport) {
				t.Errorf("WeatherAt() gotReport = %v, want %v", gotReport, test.wantReport)
			}
		})
	}
}

type mockTransport struct {
	URL string
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http"  // Override scheme to match mock server
	req.URL.Host = m.URL[7:] // Remove "http://" from mock server URL
	return http.DefaultTransport.RoundTrip(req)
}

func compareReports(a, b Report) bool {
	return a.Condition == b.Condition &&
		a.Temperature == b.Temperature &&
		a.TemperatureRaw == b.TemperatureRaw
}

// mockRoundTripper mocks the RoundTrip function for http.Client
type mockRoundTripper struct {
	fun func(req *http.Request) *http.Response
	err error
}

// RoundTrip executes the mock round trip function
func (m mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.fun(req), m.err
}

// mockHTTPClient helps in setting up a client with a mock transport
func mockHTTPClient(fn func(req *http.Request) *http.Response, err error) *http.Client {
	mrt := &mockRoundTripper{
		fun: fn,
		err: err,
	}
	return &http.Client{
		Transport: mrt,
	}
}
