package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"server/internal/model"
)

type RedisServiceInterface interface {
	GetForecast(ctx context.Context, key string) (*model.Forecast, error)
	SetForecast(ctx context.Context, key string, forecast *model.Forecast, expiration time.Duration) error
}

type RedisService struct {
	client *redis.Client
}

const (
	MAX_EXPIRATION = 24 * time.Hour
)

var ctx = context.Background()

func NewRedisService(addr, password string, db int) *RedisService {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if client.Options().Addr == "" {
		log.Fatal("CACHE_ADDRESS env variable not set")
	}

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	log.Println("Connected to Redis server")

	return &RedisService{client: client}
}

// SetForecast stores forecast data in the Redis cache.
func (r *RedisService) SetForecast(ctx context.Context, key string, forecast *model.Forecast, expiration time.Duration) error {
	if expiration <= 0 || expiration > MAX_EXPIRATION {
		return fmt.Errorf("expiration must be 0-24 hours")
	}

	json, err := json.Marshal(forecast)

	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, json, expiration).Err()
}

// GetForecast retrieves forecast data from the Redis cache.
func (r *RedisService) GetForecast(ctx context.Context, key string) (*model.Forecast, error) {

	val, err := r.client.Get(ctx, key).Result()

	if err == redis.Nil {
		return nil, fmt.Errorf("key does not exist: %s", key)
	} else if err != nil {
		return nil, fmt.Errorf("error retreiving key from Redis: %w", err)
	}

	var forecast model.Forecast
	err = json.Unmarshal([]byte(val), &forecast)

	if err != nil {
		return nil, err
	}

	return &forecast, nil
}
