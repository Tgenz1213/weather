package service_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"server/internal/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Mock census.gov response
func setupTestHandlerFunc(statusCode int, resp []byte) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write(resp)
	})
}

func TestGetCoordinates_Success(t *testing.T) {
	resp := []byte(`{
		"result": {
			"addressMatches": [{
				"coordinates": {
					"x": -118.4065,
					"y": 34.0901
				}
			}]
		}
	}`)

	handler := setupTestHandlerFunc(http.StatusOK, resp)

	server := httptest.NewServer(handler)
	defer server.Close()

	server.Client().Timeout = 10 * time.Second

	geoService := service.NewGeocodingService(server.URL, server.Client())

	ctx := context.Background()
	location, err := geoService.GetCoordinates(ctx, "123 Main St", "90210")

	assert.NoError(t, err)
	assert.NotNil(t, location)
	assert.Equal(t, 34.0901, location.Lat)
	assert.Equal(t, -118.4065, location.Lon)
}

func TestGetCoordinates_NoMatches(t *testing.T) {
	resp := []byte(`{
		"result": {
			"addressMatches": []
		}
	}`)

	handler := setupTestHandlerFunc(http.StatusOK, resp)

	server := httptest.NewServer(handler)
	defer server.Close()

	server.Client().Timeout = 10 * time.Second

	geoService := service.NewGeocodingService(server.URL, server.Client())

	ctx := context.Background()
	location, err := geoService.GetCoordinates(ctx, "123 Main St", "90210")

	assert.Error(t, err)
	assert.Nil(t, location)
	assert.Contains(t, err.Error(), "no coordinates found")
}

func TestGetCoordinates_HttpError(t *testing.T) {
	handler := setupTestHandlerFunc(http.StatusInternalServerError, []byte(nil))
	server := httptest.NewServer(handler)
	defer server.Close()

	server.Client().Timeout = 10 * time.Second

	geoService := service.NewGeocodingService(server.URL, server.Client())

	ctx := context.Background()
	location, err := geoService.GetCoordinates(ctx, "123 Main St", "90210")

	assert.Error(t, err)
	assert.Nil(t, location)
	assert.Contains(t, err.Error(), "unexpected status code")
}
