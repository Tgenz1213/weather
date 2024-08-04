package api

import (
	"net/http"
	"server/internal/controller"
)

// RegisterRoutes registers the HTTP routes for the application.
// It sets up the handler for the "/weather" endpoint, which is managed by the provided WeatherController.
//
// Parameters:
//   - mux: The HTTP request multiplexer to which the routes will be registered.
//   - weatherController: The controller responsible for handling weather-related requests.
func RegisterRoutes(mux *http.ServeMux, weatherController *controller.WeatherController) {
	mux.HandleFunc("/weather", weatherController.GetWeather)
}
