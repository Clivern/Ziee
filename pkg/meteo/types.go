// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package meteo

// City holds geocoding result for a place name.
type City struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Elevation   float64 `json:"elevation"`
	FeatureCode string  `json:"feature_code"`
	CountryCode string  `json:"country_code"`
	Admin1ID    int64   `json:"admin1_id"`
	Admin2ID    int64   `json:"admin2_id"`
	Timezone    string  `json:"timezone"`
	Population  int64   `json:"population"`
	CountryID   int64   `json:"country_id"`
	Country     string  `json:"country"`
	Admin1      string  `json:"admin1"`
	Admin2      string  `json:"admin2"`
}

type geocodingResponse struct {
	Results          []City  `json:"results"`
	GenerationTimeMS float64 `json:"generationtime_ms"`
}

// Weather holds current weather for a location.
type Weather struct {
	Latitude             float64             `json:"latitude"`
	Longitude            float64             `json:"longitude"`
	GenerationTimeMS     float64             `json:"generationtime_ms"`
	UTCOffsetSeconds     int                 `json:"utc_offset_seconds"`
	Timezone             string              `json:"timezone"`
	TimezoneAbbreviation string              `json:"timezone_abbreviation"`
	Elevation            float64             `json:"elevation"`
	CurrentWeatherUnits  CurrentWeatherUnits `json:"current_weather_units"`
	CurrentWeather       CurrentWeather      `json:"current_weather"`
}

// CurrentWeatherUnits describes units for current weather fields.
type CurrentWeatherUnits struct {
	Time          string `json:"time"`
	Interval      string `json:"interval"`
	Temperature   string `json:"temperature"`
	Windspeed     string `json:"windspeed"`
	Winddirection string `json:"winddirection"`
	IsDay         string `json:"is_day"`
	Weathercode   string `json:"weathercode"`
}

// CurrentWeather holds the latest observed weather.
type CurrentWeather struct {
	Time          string  `json:"time"`
	Interval      int     `json:"interval"`
	Temperature   float64 `json:"temperature"`
	Windspeed     float64 `json:"windspeed"`
	Winddirection int     `json:"winddirection"`
	IsDay         int     `json:"is_day"`
	Weathercode   int     `json:"weathercode"`
}
