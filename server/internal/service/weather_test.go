package service_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"server/internal/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetForecast(t *testing.T) {
	// Create a test server with a mux to handle both endpoints
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	// Mock HTTP route for /points/{latitude},{longitude}
	mux.HandleFunc("/points/34.0901,-118.4065", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/geo+json")
		w.WriteHeader(http.StatusOK)
		forecastURL := fmt.Sprintf("%s/gridpoints/LOX/149,48/forecast", server.URL)
		w.Write([]byte(fmt.Sprintf(`{
				"properties": {
					"forecast": "%s"
				}
			}`, forecastURL)))
	})

	// Mock HTTP route for /gridpoints/LOX/149,48/forecast
	mux.HandleFunc("/gridpoints/LOX/149,48/forecast", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
				"properties": {
					"periods": [{
						"number": 1,
						"name": "Today",
						"startTime": "2024-07-31T09:00:00-07:00",
						"endTime": "2024-07-31T18:00:00-07:00",
						"isDaytime": true,
						"temperature": 79,
						"temperatureUnit": "F",
						"windSpeed": "0 to 10 mph",
						"windDirection": "S",
						"shortForecast": "Sunny",
						"detailedForecast": "Sunny, with a high near 79. South wind 0 to 10 mph."
					}]
				}
			}`))
	})

	// Create WeatherService with test server URL
	weatherService := service.NewWeatherService(server.URL, &http.Client{Timeout: time.Second * 10})

	ctx := context.Background()
	forecast, err := weatherService.GetForecast(ctx, 34.0901, -118.4065)

	assert.NoError(t, err)
	assert.NotNil(t, forecast)

	expectedStartTime, _ := time.Parse(time.RFC3339, "2024-07-31T09:00:00-07:00")
	expectedEndTime, _ := time.Parse(time.RFC3339, "2024-07-31T18:00:00-07:00")
	expectedDetailedForecast := "Sunny, with a high near 79. South wind 0 to 10 mph."

	assert.Equal(t, expectedStartTime, forecast.StartTime)
	assert.Equal(t, expectedEndTime, forecast.EndTime)
	assert.Equal(t, expectedDetailedForecast, forecast.DetailedForecast)
}
