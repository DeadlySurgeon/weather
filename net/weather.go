package net

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func (s *Server) weatherAt() func(rw http.ResponseWriter, r *http.Request) {
	// Allows for configuration if need be.
	return func(rw http.ResponseWriter, r *http.Request) {
		lat := r.URL.Query().Get("lat")
		lon := r.URL.Query().Get("lon")

		if lat == "" {
			writeError(rw, http.StatusBadRequest, "invalid lat")
			return
		}

		if lon == "" {
			writeError(rw, http.StatusBadRequest, "invalid lon")
			return
		}

		report, err := s.weather.At(lat, lon)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get weather at [%v,%v]: %v\n", lat, lon, err)
			writeError(rw, http.StatusInternalServerError, "failed to get weather")
			return
		}

		e := json.NewEncoder(rw)
		e.SetIndent("", "  ")
		e.Encode(report)
	}
}

func writeError(rw http.ResponseWriter, code int, s string, i ...interface{}) {
	rw.WriteHeader(code)
	e := json.NewEncoder(rw)
	e.SetIndent("", "  ")
	e.Encode(errorResponse{Message: fmt.Sprintf(s, i...)})
}

type errorResponse struct {
	Message string `json:"message"`
}
