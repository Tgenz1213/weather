package main

import (
	"log"
	"net/http"
	"os"
	"server/api"
	"server/internal/controller"
	"server/internal/service"
	"time"
)

const (
	WEATHER_BASE_API_URL = "https://api.weather.gov"
	CENSUS_BASE_API_URL  = "https://geocoding.geo.census.gov"
)

// main is the entry point of the application. It initializes the logger, HTTP client, services, controller,
// and HTTP routes, then starts the HTTP server.
func main() {
	logger := log.New(os.Stdout, "weather_server: ", log.LstdFlags)

	mux := http.NewServeMux()

	httpClient := &http.Client{Timeout: time.Second * 10}

	gs := service.NewGeocodingService(CENSUS_BASE_API_URL, httpClient)
	ws := service.NewWeatherService(WEATHER_BASE_API_URL, httpClient)
	rs := service.NewRedisService(os.Getenv("CACHE_ADDRESS"), "", 0)
	wc := controller.NewWeatherController(gs, ws, rs)

	api.RegisterRoutes(mux, wc)

	serverAddress := os.Getenv("SERVER_ADDRESS")

	if serverAddress == "" {
		logger.Fatal("SERVER_ADDRESS environment variable is not set")
	}

	logger.Printf("Starting server on %s", serverAddress)

	if err := http.ListenAndServe(serverAddress, mux); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}
