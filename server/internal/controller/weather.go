package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/internal/model"
	"server/internal/service"
	"time"
)

type WeatherController struct {
	geocodingService service.GeocodingServiceInterface
	weatherService   service.WeatherServiceInterface
	redisService     service.RedisServiceInterface
}

func NewWeatherController(gs service.GeocodingServiceInterface, ws service.WeatherServiceInterface, rs service.RedisServiceInterface) *WeatherController {
	return &WeatherController{
		geocodingService: gs,
		weatherService:   ws,
		redisService:     rs,
	}
}

// GetWeather handles the HTTP request to retrieve weather data.
func (wc *WeatherController) GetWeather(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	street := r.URL.Query().Get("street")
	zip := r.URL.Query().Get("zip")

	location, err := wc.getCoordinates(ctx, street, zip)

	if err != nil {
		wc.handleError(w, err, http.StatusInternalServerError)
		return
	}

	cacheKey := "forecast_" + zip

	forecast, err := wc.getForecastFromApi(ctx, location)

	if err != nil {
		wc.handleError(w, err, http.StatusInternalServerError)
		return
	}

	if err = wc.cacheForecast(ctx, cacheKey, forecast); err != nil {
		wc.handleError(w, err, http.StatusInternalServerError)
		return
	}

	if forecast, err = wc.getForecastFromCache(ctx, cacheKey); err != nil {
		wc.handleError(w, err, http.StatusInternalServerError)
		return
	}

	wc.respondWithJSON(w, forecast)
}

// getCoordinates retrieves coordinates from the geocoding service.
func (wc *WeatherController) getCoordinates(ctx context.Context, street, zipCode string) (*service.Location, error) {
	location, err := wc.geocodingService.GetCoordinates(ctx, street, zipCode)

	if err != nil {
		return nil, fmt.Errorf("error getting coordinates: %w", err)
	}

	return location, nil
}

// getForecastFromApi retrieves forecast data from the weather service.
func (wc *WeatherController) getForecastFromApi(ctx context.Context, location *service.Location) (*model.Forecast, error) {
	forecast, err := wc.weatherService.GetForecast(ctx, location.Lat, location.Lon)

	if err != nil {
		return nil, fmt.Errorf("error fetching forecast: %w", err)
	}

	return forecast, nil
}

// cacheForecast stores forecast data in the Redis cache.
func (wc *WeatherController) cacheForecast(ctx context.Context, cacheKey string, forecast *model.Forecast) error {
	err := wc.redisService.SetForecast(ctx, cacheKey, forecast, 1*time.Hour)

	if err != nil {
		return fmt.Errorf("error caching forecast: %w", err)
	}

	return nil
}

// getForecastFromCache retrieves forecast data from the Redis cache.
func (wc *WeatherController) getForecastFromCache(ctx context.Context, cacheKey string) (*model.Forecast, error) {
	forecast, err := wc.redisService.GetForecast(ctx, cacheKey)

	if err != nil {
		return nil, fmt.Errorf("error retreiving forecast: %w", err)
	}

	return forecast, nil
}

// handleError logs and returns errors.
func (wc *WeatherController) handleError(w http.ResponseWriter, err error, statusCode int) {
	log.Printf("Error: %v", err)
	http.Error(w, err.Error(), statusCode)
}

// respondWithJSON sends a JSON response.
func (wc *WeatherController) respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		wc.handleError(w, err, http.StatusInternalServerError)
	}
}
