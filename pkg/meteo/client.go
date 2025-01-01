// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package meteo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
)

const (
	defaultGeocodingURL = "https://geocoding-api.open-meteo.com/v1"
	defaultForecastURL  = "https://api.open-meteo.com/v1"
)

// Client fetches city and weather data from Open-Meteo.
type Client interface {
	GetCity(ctx context.Context, name string, language string) (*City, error)
	GetWeatherByCity(ctx context.Context, name string, language string) (*Weather, error)
	GetWeatherByCoordinates(ctx context.Context, latitude, longitude float64) (*Weather, error)
}

// openMeteoClient calls the Open-Meteo geocoding and forecast APIs.
type openMeteoClient struct {
	httpClient   *http.Client
	geocodingURL string
	forecastURL  string
}

// Option configures a Client.
type Option func(*openMeteoClient)

// WithHTTPClient sets the HTTP client used for API requests.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *openMeteoClient) {
		c.httpClient = httpClient
	}
}

// WithGeocodingURL overrides the geocoding API base URL.
func WithGeocodingURL(url string) Option {
	return func(c *openMeteoClient) {
		c.geocodingURL = url
	}
}

// WithForecastURL overrides the forecast API base URL.
func WithForecastURL(url string) Option {
	return func(c *openMeteoClient) {
		c.forecastURL = url
	}
}

// New returns a Client with sensible defaults.
func New(opts ...Option) Client {
	c := &openMeteoClient{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		geocodingURL: defaultGeocodingURL,
		forecastURL:  defaultForecastURL,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// GetCity resolves a place name to geographic details using the Open-Meteo geocoding API.
func (c *openMeteoClient) GetCity(ctx context.Context, name string, language string) (*City, error) {
	name = strings.TrimSpace(name)
	if lo.IsEmpty(name) {
		return nil, ErrEmptyCityName
	}

	if lo.IsEmpty(language) {
		language = "en"
	}

	query := url.Values{
		"name":     {name},
		"count":    {"1"},
		"language": {language},
		"format":   {"json"},
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.geocodingURL+"/search?"+query.Encode(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("meteo geocoding request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("meteo geocoding request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meteo geocoding request: unexpected status %s", resp.Status)
	}

	var payload geocodingResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("meteo geocoding decode: %w", err)
	}

	if len(payload.Results) == 0 {
		return nil, ErrCityNotFound
	}

	return &payload.Results[0], nil
}

// GetWeatherByCity resolves a city name and returns its current weather.
func (c *openMeteoClient) GetWeatherByCity(ctx context.Context, name string, language string) (*Weather, error) {
	city, err := c.GetCity(ctx, name, language)
	if err != nil {
		return nil, err
	}

	return c.GetWeatherByCoordinates(ctx, city.Latitude, city.Longitude)
}

// GetWeatherByCoordinates returns current weather for the given coordinates.
func (c *openMeteoClient) GetWeatherByCoordinates(ctx context.Context, latitude, longitude float64) (*Weather, error) {
	query := url.Values{
		"latitude":        {strconv.FormatFloat(latitude, 'f', -1, 64)},
		"longitude":       {strconv.FormatFloat(longitude, 'f', -1, 64)},
		"current_weather": {"true"},
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.forecastURL+"/forecast?"+query.Encode(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("meteo forecast request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("meteo forecast request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meteo forecast request: unexpected status %s", resp.Status)
	}

	var weather Weather
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, fmt.Errorf("meteo forecast decode: %w", err)
	}

	if lo.IsEmpty(weather.CurrentWeather.Time) {
		return nil, ErrInvalidResponse
	}

	return &weather, nil
}
