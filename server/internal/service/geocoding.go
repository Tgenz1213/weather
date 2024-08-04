// Service for interacting with Census.gov's geocoding API
// https://geocoding.geo.census.gov/geocoder/Geocoding_Services_API.pdf

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"gopkg.in/validator.v2"
)

type GeocodingServiceInterface interface {
	GetCoordinates(ctx context.Context, street, zipCode string) (*Location, error)
}

type GeocodingService struct {
	httpClient *http.Client
	baseApiURL string
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func NewGeocodingService(url string, client *http.Client) *GeocodingService {
	return &GeocodingService{
		httpClient: client,
		baseApiURL: url,
	}
}

// GetCoordinates retrieves coordinates for a given location from Census.gov's forward geocoding API.
func (g *GeocodingService) GetCoordinates(ctx context.Context, street, zip string) (*Location, error) {
	preparedStreet, err := prepareStreet(street)

	if err != nil {
		return nil, err
	}

	preparedZip, err := prepareZip(zip)

	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/geocoder/locations/address?street=%s&zip=%s&benchmark=4&format=json", g.baseApiURL, preparedStreet, preparedZip)

	log.Printf("GeocodingService.GetCoordinates URL: %s", url)

	geoData, err := g.makeRequest(ctx, url)

	if err != nil {
		return nil, err
	}

	if err := validator.Validate(geoData); err != nil {
		return nil, fmt.Errorf("invalid response received: %w", err)
	}

	if len(geoData.Result.AddressMatches) == 0 {
		return nil, fmt.Errorf("no coordinates found for street %s and zip code %s", street, zip)
	}

	return &Location{
		Lat: geoData.Result.AddressMatches[0].Coordinates.Y,
		Lon: geoData.Result.AddressMatches[0].Coordinates.X,
	}, nil
}

// makeRequest performs an HTTP GET request to the geocoding API and decodes the response.
func (g *GeocodingService) makeRequest(ctx context.Context, url string) (*GeoData, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := g.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error making request to geocoding API: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data GeoData

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("error decoding geocoding response: %w", err)
	}

	return &data, nil
}

type GeoData struct {
	Result struct {
		AddressMatches []struct {
			Coordinates struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
			} `json:"coordinates"`
		} `json:"addressMatches"`
	} `json:"result" validate:"nonzero"`
}

// prepareStreet validates and prepares a street address for use in a URL.
func prepareStreet(street string) (string, error) {
	if street == "" {
		return "", fmt.Errorf("missing street")
	}

	if len(street) > 100 {
		return "", fmt.Errorf("street is too long")
	}

	return url.QueryEscape(street), nil
}

// prepareZip validates and prepares a ZIP code for use in a URL.
func prepareZip(zip string) (string, error) {
	if len(zip) != 5 {
		return "", fmt.Errorf("zip code must be 5 digits")
	}

	if _, err := strconv.Atoi(zip); err != nil {
		return "", fmt.Errorf("zip code must be numbers only")
	}

	return url.QueryEscape(zip), nil
}
