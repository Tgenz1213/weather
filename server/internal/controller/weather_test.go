package controller_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"server/internal/controller"
	"server/internal/model"
	"server/internal/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock services
type MockGeocodingService struct {
	mock.Mock
}

func (m *MockGeocodingService) GetCoordinates(ctx context.Context, street, zipCode string) (*service.Location, error) {
	args := m.Called(ctx, street, zipCode)
	if args.Get(0) != nil {
		return args.Get(0).(*service.Location), args.Error(1)
	}
	return nil, args.Error(1)
}

type MockWeatherService struct {
	mock.Mock
}

func (m *MockWeatherService) GetForecast(ctx context.Context, lat, lon float64) (*model.Forecast, error) {
	args := m.Called(ctx, lat, lon)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Forecast), args.Error(1)
	}
	return nil, args.Error(1)
}

type MockRedisService struct {
	mock.Mock
}

func (m *MockRedisService) GetForecast(ctx context.Context, key string) (*model.Forecast, error) {
	args := m.Called(ctx, key)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Forecast), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRedisService) SetForecast(ctx context.Context, key string, forecast *model.Forecast, expiration time.Duration) error {
	args := m.Called(ctx, key, forecast, expiration)
	return args.Error(0)
}

func TestGetWeather(t *testing.T) {
	geocodingService := new(MockGeocodingService)
	weatherService := new(MockWeatherService)
	redisService := new(MockRedisService)

	wc := controller.NewWeatherController(geocodingService, weatherService, redisService)

	// Mock data
	zipCode := "90210"
	street := "123 Beverly Hills"
	cacheKey := "forecast_" + zipCode
	location := &service.Location{Lat: 34.0901, Lon: -118.4065}
	startTime, _ := time.Parse(time.RFC3339, "2024-07-31T09:00:00-07:00")
	endTime, _ := time.Parse(time.RFC3339, "2024-07-31T18:00:00-07:00")
	forecast := &model.Forecast{StartTime: startTime, EndTime: endTime, DetailedForecast: "Sunny, with a high near 79. South wind 0 to 10 mph."}

	// Mocking the behavior
	redisService.On("GetForecast", mock.Anything, cacheKey).Return(nil, nil).Once()
	geocodingService.On("GetCoordinates", mock.Anything, street, zipCode).Return(location, nil).Once()
	weatherService.On("GetForecast", mock.Anything, location.Lat, location.Lon).Return(forecast, nil).Once()
	redisService.On("SetForecast", mock.Anything, cacheKey, forecast, 1*time.Hour).Return(nil).Once()
	redisService.On("GetForecast", mock.Anything, cacheKey).Return(forecast, nil).Once()

	req, err := http.NewRequest("GET", "/weather?zip=90210&street=123+Beverly+Hills", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(wc.GetWeather)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expected := `{"startTime":"2024-07-31T09:00:00-07:00","endTime":"2024-07-31T18:00:00-07:00","detailedForecast":"Sunny, with a high near 79. South wind 0 to 10 mph."}`
	assert.JSONEq(t, expected, rr.Body.String())

	redisService.AssertExpectations(t)
	geocodingService.AssertExpectations(t)
	weatherService.AssertExpectations(t)
}
