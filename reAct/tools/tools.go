package tools

import (
	"math/rand"
	"react/models"
	"strings"
	"time"
)

var (
	GetCurrentWeatherAllProperties = map[string]any{
		"address": map[string]string{
			"type": "string",
		},
	}

	GetCurrentWeatherRequiredProperties = []string{"address"}
)

var (
	GetLocationAllProperties      = map[string]any{}
	GetLocationRequiredProperties = []string{}
)

type Tooler interface {
	GetCurrentWeather(models.Location) models.Weather
	GetLocation() models.Location
}

type HardCodedTool struct{}

func GetNewHardCodedTool() Tooler {
	return &HardCodedTool{}
}

func (t *HardCodedTool) GetCurrentWeather(loc models.Location) models.Weather {
	switch {
	case strings.Contains(strings.ToLower(loc.Address), "bhubaneswar"):
		return models.Weather{
			Temperature: "35",
			Unit:        "C",
			Forecast:    "Sunny",
		}

	case strings.Contains(strings.ToLower(loc.Address), "oslo"):
		return models.Weather{
			Temperature: "23",
			Unit:        "C",
			Forecast:    "Rainy",
		}
	}

	return models.Weather{
		Temperature: "15",
		Unit:        "C",
		Forecast:    "Windy",
	}
}

var locations = []models.Location{
	{Address: "oslo"},
	{Address: "Delta Square, Bhubaneswar, Odisha, India"},
}

func (t *HardCodedTool) GetLocation() models.Location {
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(2)

	if randomNum >= 0 && randomNum < len(locations) {
		return locations[randomNum]
	}

	return models.Location{Address: "Texas"}
}
