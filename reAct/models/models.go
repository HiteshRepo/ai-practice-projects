package models

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Weather struct {
	Temperature string `json:"temperature"`
	Unit        string `json:"unit"`
	Forecast    string `json:"forecast"`
}

func (w Weather) ToString() string {
	b, _ := json.Marshal(w)
	return string(b)
}

type Location struct {
	Address string `json:"address"`
}

func (l Location) ToString() string {
	return l.Address
}

type Action struct {
	FunctionName string
	Arguments    []string
	ToolCallID   string
}

type WeatherInputs struct {
	Address string `json:"address"`
}

type LocationInfo struct {
	City          string  `json:"city"`
	Region        string  `json:"region"`
	Country       string  `json:"country"`
	CountryName   string  `json:"country_name"`
	CountryCode   string  `json:"country_code"`
	ContinentCode string  `json:"continent_code"`
	Postal        string  `json:"postal"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
}

func (li *LocationInfo) GenerateAddress() string {
	return fmt.Sprintf("%s, %s, %s", li.City, li.Region, li.Country)
}

// https://api.weatherstack.com/current
type WeatherResponse struct {
	Request Request `json:"request"`
	Current Current `json:"current"`
}

type Request struct {
	Type     string `json:"type"`
	Query    string `json:"query"`
	Language string `json:"language"`
	Unit     string `json:"unit"`
}

type Current struct {
	ObservationTime     string     `json:"observation_time"`
	Temperature         int        `json:"temperature"`
	WeatherCode         int        `json:"weather_code"`
	WeatherIcons        []string   `json:"weather_icons"`
	WeatherDescriptions []string   `json:"weather_descriptions"`
	Astro               Astro      `json:"astro"`
	AirQuality          AirQuality `json:"air_quality"`
	WindSpeed           int        `json:"wind_speed"`
	WindDegree          int        `json:"wind_degree"`
	WindDir             string     `json:"wind_dir"`
	Pressure            int        `json:"pressure"`
	Precip              float64    `json:"precip"`
	Humidity            int        `json:"humidity"`
	CloudCover          int        `json:"cloudcover"`
	FeelsLike           int        `json:"feelslike"`
	UVIndex             int        `json:"uv_index"`
	Visibility          int        `json:"visibility"`
	IsDay               string     `json:"is_day"`
}

type Astro struct {
	Sunrise          string `json:"sunrise"`
	Sunset           string `json:"sunset"`
	Moonrise         string `json:"moonrise"`
	Moonset          string `json:"moonset"`
	MoonPhase        string `json:"moon_phase"`
	MoonIllumination int    `json:"moon_illumination"`
}

type AirQuality struct {
	CO           string `json:"co"`
	NO2          string `json:"no2"`
	O3           string `json:"o3"`
	SO2          string `json:"so2"`
	PM25         string `json:"pm2_5"`
	PM10         string `json:"pm10"`
	USEPAIndex   string `json:"us-epa-index"`
	GBDefraIndex string `json:"gb-defra-index"`
}

func ConvertWeatherResponseToWeather(wr WeatherResponse) *Weather {
	temperature := fmt.Sprintf("%d (feels like %d)", wr.Current.Temperature, wr.Current.FeelsLike)

	unit := strings.ToUpper(wr.Request.Unit)
	if unit == "M" {
		unit = "Metric (Celsius)"
	} else if unit == "F" {
		unit = "Imperial (Fahrenheit)"
	}

	var forecastParts []string

	// Current weather condition
	if len(wr.Current.WeatherDescriptions) > 0 {
		forecastParts = append(forecastParts, fmt.Sprintf("🌤️ Condition: %s",
			strings.Join(wr.Current.WeatherDescriptions, ", ")))
	}

	// Atmospheric conditions
	forecastParts = append(forecastParts, fmt.Sprintf("💨 Wind: %d km/h %s (%d°)",
		wr.Current.WindSpeed, wr.Current.WindDir, wr.Current.WindDegree))
	forecastParts = append(forecastParts, fmt.Sprintf("💧 Humidity: %d%% | Pressure: %d mbar",
		wr.Current.Humidity, wr.Current.Pressure))
	forecastParts = append(forecastParts, fmt.Sprintf("🌧️ Precipitation: %.1f mm | Cloud Cover: %d%%",
		wr.Current.Precip, wr.Current.CloudCover))
	forecastParts = append(forecastParts, fmt.Sprintf("👁️ Visibility: %d km | UV Index: %d",
		wr.Current.Visibility, wr.Current.UVIndex))

	// Day/Night status
	dayNight := "Day"
	if wr.Current.IsDay == "no" {
		dayNight = "Night"
	}
	forecastParts = append(forecastParts, fmt.Sprintf("🌅 Time of Day: %s", dayNight))

	// Astronomical information
	forecastParts = append(forecastParts, fmt.Sprintf("🌅 Sunrise: %s | Sunset: %s",
		wr.Current.Astro.Sunrise, wr.Current.Astro.Sunset))
	forecastParts = append(forecastParts, fmt.Sprintf("🌙 Moon: %s (%d%% illuminated)",
		wr.Current.Astro.MoonPhase, wr.Current.Astro.MoonIllumination))
	forecastParts = append(forecastParts, fmt.Sprintf("🌛 Moonrise: %s | Moonset: %s",
		wr.Current.Astro.Moonrise, wr.Current.Astro.Moonset))

	// Air quality information
	forecastParts = append(forecastParts, fmt.Sprintf("🏭 Air Quality (US EPA Index: %s, UK DEFRA: %s)",
		wr.Current.AirQuality.USEPAIndex, wr.Current.AirQuality.GBDefraIndex))
	forecastParts = append(forecastParts, fmt.Sprintf("🌫️ Pollutants - PM2.5: %s μg/m³, PM10: %s μg/m³",
		wr.Current.AirQuality.PM25, wr.Current.AirQuality.PM10))
	forecastParts = append(forecastParts, fmt.Sprintf("☁️ Gases - CO: %s, NO2: %s, O3: %s, SO2: %s",
		wr.Current.AirQuality.CO, wr.Current.AirQuality.NO2, wr.Current.AirQuality.O3, wr.Current.AirQuality.SO2))

	forecast := strings.Join(forecastParts, "\n")

	return &Weather{
		Temperature: temperature,
		Unit:        unit,
		Forecast:    forecast,
	}
}
