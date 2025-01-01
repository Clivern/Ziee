// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"context"

	"github.com/actx0/ziee/pkg/meteo"
)

// Tools runs external tool integrations.
type Tools struct {
	Meteo meteo.Client
}

// NewTools creates a tools module.
func NewTools(meteo meteo.Client) *Tools {
	return &Tools{
		Meteo: meteo,
	}
}

// GetCityInformationRequest is the body for the getCityInformation tool.
type GetCityInformationRequest struct {
	City     string `json:"city" validate:"required" label:"City"`
	Language string `json:"language" validate:"omitempty" label:"Language"`
}

// GetWeatherByCityRequest is the body for the getWeatherByCity tool.
type GetWeatherByCityRequest struct {
	City     string `json:"city" validate:"required" label:"City"`
	Language string `json:"language" validate:"omitempty" label:"Language"`
}

// GetWeatherByCoordinatesRequest is the body for the getWeatherByCoordinates tool.
type GetWeatherByCoordinatesRequest struct {
	Latitude  float64 `json:"latitude" validate:"required" label:"Latitude"`
	Longitude float64 `json:"longitude" validate:"required" label:"Longitude"`
}

// GetCityInformation resolves a city name to geographic details.
func (t *Tools) GetCityInformation(ctx context.Context, req *GetCityInformationRequest) (*meteo.City, error) {
	return t.Meteo.GetCity(ctx, req.City, req.Language)
}

// GetWeatherByCity resolves a city name and returns its current weather.
func (t *Tools) GetWeatherByCity(ctx context.Context, req *GetWeatherByCityRequest) (*meteo.Weather, error) {
	return t.Meteo.GetWeatherByCity(ctx, req.City, req.Language)
}

// GetWeatherByCoordinates returns current weather for the given coordinates.
func (t *Tools) GetWeatherByCoordinates(ctx context.Context, req *GetWeatherByCoordinatesRequest) (*meteo.Weather, error) {
	return t.Meteo.GetWeatherByCoordinates(ctx, req.Latitude, req.Longitude)
}
