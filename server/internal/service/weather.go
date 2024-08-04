package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"server/internal/model"
	"time"
)

type WeatherServiceInterface interface {
	GetForecast(ctx context.Context, lat, lon float64) (*model.Forecast, error)
}

type WeatherService struct {
	apiBaseURL string
	httpClient *http.Client
	headers    http.Header
}

func NewWeatherService(url string, client *http.Client) *WeatherService {
	headers := make(http.Header)
	headers.Set("User-Agent", os.Getenv("USER_AGENT"))
	headers.Set("Accept", "application/geo+json")

	return &WeatherService{
		apiBaseURL: url,
		httpClient: client,
		headers:    headers,
	}
}

// GetForecast retrieves the weather forecast for the given latitude and longitude using the weather.gov API.
func (w *WeatherService) GetForecast(ctx context.Context, latitude, longitude float64) (*model.Forecast, error) {
	pointURL := fmt.Sprintf("%s/points/%.4f,%.4f", w.apiBaseURL, latitude, longitude)

	// Get point data for coordinates
	pointData, err := w.makeRequest(ctx, pointURL)

	if err != nil {
		return nil, fmt.Errorf("error fetching point data: %w", err)
	}

	// Get forecast data from forecast URL in point data
	forecastData, err := w.makeRequest(ctx, pointData.Properties.Forecast)

	if err != nil {
		return nil, fmt.Errorf("error fetching forecast data: %w", err)
	}

	if len(forecastData.Properties.Periods) == 0 {
		return nil, fmt.Errorf("no forecast data available")
	}

	currentConditions := forecastData.Properties.Periods[0]

	return &model.Forecast{
		StartTime:        currentConditions.StartTime,
		EndTime:          currentConditions.EndTime,
		DetailedForecast: currentConditions.DetailedForecast,
	}, nil
}

// makeRequest performs an HTTP GET request and decodes the response.
func (w *WeatherService) makeRequest(ctx context.Context, url string) (*WeatherData, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header = w.headers.Clone()

	resp, err := w.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data WeatherData

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &data, nil
}

type WeatherData struct {
	Properties struct {
		Forecast string     `json:"forecast"` // Point URL response body
		Periods  []struct { // Forecast response body
			StartTime        time.Time `json:"startTime"`
			EndTime          time.Time `json:"endTime"`
			DetailedForecast string    `json:"detailedForecast"`
		} `json:"periods"`
	} `json:"properties"`
}
