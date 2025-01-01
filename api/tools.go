// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package api

import (
	"errors"
	"net/http"

	"github.com/actx0/ziee/locale"
	"github.com/actx0/ziee/module"
	"github.com/actx0/ziee/pkg/meteo"
	"github.com/actx0/ziee/pkg/util"

	"github.com/rs/zerolog/log"
)

// GetCityInformationAction runs the getCityInformation tool.
func GetCityInformationAction(w http.ResponseWriter, r *http.Request) {
	var req module.GetCityInformationRequest
	err := util.DecodeAndValidate(r, &req)
	if err != nil {
		util.WriteValidationError(w, err)
		return
	}

	tools := module.NewTools(meteo.New())

	city, err := tools.GetCityInformation(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, meteo.ErrCityNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "city_not_found"),
			})
		case errors.Is(err, meteo.ErrEmptyCityName):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "city_required"),
			})
		default:
			log.Error().
				Err(err).
				Str("toolName", "getCityInformation").
				Msg("Failed to run tool")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_run_tool"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, city)
}

// GetWeatherByCityAction runs the getWeatherByCity tool.
func GetWeatherByCityAction(w http.ResponseWriter, r *http.Request) {
	var req module.GetWeatherByCityRequest
	err := util.DecodeAndValidate(r, &req)
	if err != nil {
		util.WriteValidationError(w, err)
		return
	}

	tools := module.NewTools(meteo.New())

	weather, err := tools.GetWeatherByCity(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, meteo.ErrCityNotFound):
			util.WriteJSON(w, http.StatusNotFound, map[string]any{
				"errorMessage": locale.TR(r, "city_not_found"),
			})
		case errors.Is(err, meteo.ErrEmptyCityName):
			util.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errorMessage": locale.TR(r, "city_required"),
			})
		default:
			log.Error().
				Err(err).
				Str("toolName", "getWeatherByCity").
				Msg("Failed to run tool")
			util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"errorMessage": locale.TR(r, "failed_run_tool"),
			})
		}
		return
	}

	util.WriteJSON(w, http.StatusOK, weather)
}

// GetWeatherByCoordinatesAction runs the getWeatherByCoordinates tool.
func GetWeatherByCoordinatesAction(w http.ResponseWriter, r *http.Request) {
	var req module.GetWeatherByCoordinatesRequest
	err := util.DecodeAndValidate(r, &req)
	if err != nil {
		util.WriteValidationError(w, err)
		return
	}

	tools := module.NewTools(meteo.New())

	weather, err := tools.GetWeatherByCoordinates(r.Context(), &req)
	if err != nil {
		log.Error().
			Err(err).
			Str("toolName", "getWeatherByCoordinates").
			Msg("Failed to run tool")
		util.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"errorMessage": locale.TR(r, "failed_run_tool"),
		})
		return
	}

	util.WriteJSON(w, http.StatusOK, weather)
}
