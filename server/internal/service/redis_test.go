package service_test

import (
	"context"
	"server/internal/model"
	"server/internal/service"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

type forecastTestCase struct {
	name           string
	key            string
	expiration     time.Duration
	fastForward    time.Duration
	forecast       *model.Forecast
	expectedResult *model.Forecast
	expectedError  error
	setErrMsg      string
	getErrMsg      string
	getReturnMsg   string
}

var forecastTestCases = []forecastTestCase{
	{
		name:           "Set and get success",
		key:            "test_forecast",
		expiration:     10 * time.Second,
		fastForward:    0 * time.Second,
		forecast:       &model.Forecast{DetailedForecast: "Sunny, with a high near 79. South wind 0 to 10 mph."},
		expectedResult: &model.Forecast{DetailedForecast: "Sunny, with a high near 79. South wind 0 to 10 mph."},
		expectedError:  nil,
		setErrMsg:      "Setting forecast should not return an error",
		getErrMsg:      "Getting forecast should not return an error",
		getReturnMsg:   "Getting forecast should return same forecast",
	},
	{
		name:           "",
		key:            "",
		expiration:     0 * time.Second,
		fastForward:    0 * time.Second,
		forecast:       &model.Forecast{DetailedForecast: "Sunny, with a high near 79. South wind 0 to 10 mph."},
		expectedResult: &model.Forecast{DetailedForecast: "Sunny, with a high near 79. South wind 0 to 10 mph."},
		expectedError:  nil,
	},
}

func parseTime(t *testing.T, value string) time.Time {
	parsedTime, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatalf("time '%s' could not be parsed", err)
	}
	return parsedTime
}

func setupTestRedis(t *testing.T) (*service.RedisService, *miniredis.Miniredis, func()) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	redisService := service.NewRedisService(mr.Addr(), "", 0)
	return redisService, mr, func() {
		mr.Close()
	}
}

func TestSetAndGetForecast(t *testing.T) {
	redisService, _, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()

	key := "test_forecast"
	startTime, _ := time.Parse(time.RFC3339, "2024-07-31T09:00:00-07:00")
	endTime, _ := time.Parse(time.RFC3339, "2024-07-31T18:00:00-07:00")
	forecast := &model.Forecast{
		StartTime:        startTime,
		EndTime:          endTime,
		DetailedForecast: "Sunny, with a high near 79. South wind 0 to 10 mph.",
	}

	err := redisService.SetForecast(ctx, key, forecast, time.Hour)
	assert.NoError(t, err, "Setting forecast should not return an error")

	retrievedForecast, err := redisService.GetForecast(ctx, key)
	if assert.NoError(t, err, "getting forecast should not return an error") {
		assert.Equal(t, forecast, retrievedForecast)
	}
}

func TestGetNonExistentForecast(t *testing.T) {
	redisService, _, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	key := "non_existent_forecast"

	forecast, err := redisService.GetForecast(ctx, key)
	if assert.NoError(t, err, "Getting nonexistent key should not return an error") {
		assert.Nil(t, forecast, "Getting nonexistent key should return nil forecast")
	}
}

func TestExpiration(t *testing.T) {
	redisService, mr, cleanup := setupTestRedis(t)
	defer cleanup()

	const TIMEOUT = 10

	ctx := context.Background()
	key := "expiring_forecast"
	forecast := &model.Forecast{DetailedForecast: "Test"}

	err := redisService.SetForecast(ctx, key, forecast, TIMEOUT*time.Second)
	assert.NoError(t, err, "Setting forecast should not return an error")

	retrievedForecast, err := redisService.GetForecast(ctx, key)
	if assert.NoError(t, err, "Getting non-expired forecast should not produce an error") {
		assert.Equal(t, forecast, retrievedForecast, "Retrieved forecast should match set forecast")
	}

	mr.FastForward((TIMEOUT + 1) * time.Second)

	forecast, err = redisService.GetForecast(ctx, key)
	if assert.NoError(t, err, "Getting expired forecast should not produce an error") {
		assert.Nil(t, forecast, "Getting expired forecast should return nil forecast")
	}
}

func TestSetForecastWithInvalidExpiration(t *testing.T) {
	redisService, _, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	key := "invalid_expiration_forecast"
	forecast := &model.Forecast{DetailedForecast: "Test"}

	err := redisService.SetForecast(ctx, key, forecast, 25*time.Hour)
	assert.Error(t, err, "Setting expiration above 24 hours should return an error")

	err = redisService.SetForecast(ctx, key, forecast, 0)
	assert.Error(t, err, "Setting 0 expiration or lower should return an error")
}
