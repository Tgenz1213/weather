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

	// Test SetForecast
	err := redisService.SetForecast(ctx, key, forecast, time.Hour)
	assert.NoError(t, err, "TestSetAndGetForecast: SetForecast assert no error")

	// Test GetForecast
	retrievedForecast, err := redisService.GetForecast(ctx, key)
	assert.NoError(t, err, "TestSetAndGetForecast: GetForecast assert no error")
	assert.Equal(t, forecast, retrievedForecast)
}

func TestGetNonExistentForecast(t *testing.T) {
	redisService, _, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	key := "non_existent_forecast"

	forecast, err := redisService.GetForecast(ctx, key)
	assert.NoError(t, err, "Getting nonexistent key should not return an error")
	assert.Nil(t, forecast, "Getting nonexistent key should return nil forecast")
}

func TestExpiration(t *testing.T) {
	redisService, mr, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	key := "expiring_forecast"
	forecast := &model.Forecast{DetailedForecast: "Test"}

	t.Run("Key expires after set time", func(t *testing.T) {
		err := redisService.SetForecast(ctx, key, forecast, 1*time.Second)
		assert.NoError(t, err, "Setting forecast should not produce an error")

		mr.FastForward(2 * time.Second)

		forecast, err = redisService.GetForecast(ctx, key)
		assert.NoError(t, err, "Getting expired forecast should not produce an error")
		assert.Nil(t, forecast, "Getting expired forecast should return nil forecast")
	})

	t.Run("Key does not expire before set time", func(t *testing.T) {
		err := redisService.SetForecast(ctx, key, forecast, 3*time.Second)
		assert.NoError(t, err, "Setting forecast should not produce an error")

		mr.FastForward(1 * time.Second)

		retrievedForecast, err := redisService.GetForecast(ctx, key)
		assert.NoError(t, err, "Getting non-expired forecast should not produce an error")
		assert.Equal(t, forecast, retrievedForecast, "Retrieved forecast should match set forecast")
	})
}

func TestSetForecastWithInvalidExpiration(t *testing.T) {
	redisService, _, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	key := "invalid_expiration_forecast"
	forecast := &model.Forecast{DetailedForecast: "Test"}

	err := redisService.SetForecast(ctx, key, forecast, 25*time.Hour)
	assert.Error(t, err, "TestSetForecastWithInavlidExpiration: GetForecast assert error expiration above max")

	err = redisService.SetForecast(ctx, key, forecast, 0)
	assert.Error(t, err, "TestSetForecastWithInavlidExpiration: GetForecast assert error expiration below min")
}
