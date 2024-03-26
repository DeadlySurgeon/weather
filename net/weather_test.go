package net

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/deadlysurgeon/weather/weather"
	"github.com/deadlysurgeon/weather/weather/mock"
	"go.uber.org/mock/gomock"
)

func TestWeatherAtHandler(t *testing.T) {
	for name, test := range map[string]struct {
		name           string
		lat            string
		lon            string
		mockReturn     interface{} // Use interface{} to handle both WeatherReport and error
		mockReturnErr  error
		expectedCode   int
		expectedResult string
	}{
		"valid request": {
			lat:            "40.7128",
			lon:            "-74.0060",
			mockReturn:     weather.Report{Condition: "snow", Temperature: "freezing", TemperatureRaw: -2.7},
			mockReturnErr:  nil,
			expectedCode:   http.StatusOK,
			expectedResult: `"condition": "snow"`,
		},
		"mock return error": {
			lat:            "40.7128",
			lon:            "-74.0060",
			mockReturn:     weather.Report{Condition: "snow", Temperature: "freezing", TemperatureRaw: -2.7},
			mockReturnErr:  fmt.Errorf("failed to do thing"),
			expectedCode:   http.StatusInternalServerError,
			expectedResult: `"message": "failed to get weather"`,
		},
		"missing latitude": {

			lat:            "",
			lon:            "-74.0060",
			expectedCode:   http.StatusBadRequest,
			expectedResult: `"message": "invalid lat"`,
		},
		"missing longitude": {
			lat:            "40.7128",
			lon:            "",
			expectedCode:   http.StatusBadRequest,
			expectedResult: `"message": "invalid lon"`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockWeatherService := mock.NewMockService(ctrl)
			if test.lat != "" && test.lon != "" { // Set up mock only for valid lat and lon
				mockWeatherService.EXPECT().
					At(test.lat, test.lon).
					Return(test.mockReturn, test.mockReturnErr)
			}

			s := &Server{weather: mockWeatherService}

			req, err := http.NewRequest("GET", "example/", nil)
			if err != nil {
				t.Fatal(err)
			}

			q := req.URL.Query()
			q.Add("lat", test.lat)
			q.Add("lon", test.lon)
			req.URL.RawQuery = q.Encode()

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(s.weatherAt())

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != test.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, test.expectedCode)
			}

			if !bytes.Contains(rr.Body.Bytes(), []byte(test.expectedResult)) {
				t.Errorf("handler returned unexpected body: expected to contain %v",
					test.expectedResult)
			}
		})
	}
}
