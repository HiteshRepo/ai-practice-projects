package tools

import (
	"math/rand"
	"react/models"
	"react/utils"
	"strings"
	"time"
)

type HardCodedTool struct{}

func GetNewHardCodedTool() Tooler {
	return &HardCodedTool{}
}

func (t *HardCodedTool) GetCurrentWeather(loc models.Location) (models.Weather, error) {
	switch {
	case strings.Contains(strings.ToLower(loc.Address), "bhubaneswar"):
		return models.Weather{
			Temperature: "35",
			Unit:        "C",
			Forecast:    "Sunny",
		}, nil

	case strings.Contains(strings.ToLower(loc.Address), "oslo"):
		return models.Weather{
			Temperature: "23",
			Unit:        "C",
			Forecast:    "Rainy",
		}, nil
	}

	return models.Weather{
		Temperature: "15",
		Unit:        "C",
		Forecast:    "Windy",
	}, nil
}

var locations = []models.Location{
	{Address: "oslo"},
	{Address: "Delta Square, Bhubaneswar, Odisha, India"},
}

func (t *HardCodedTool) GetLocation() (models.Location, error) {
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(2)

	if randomNum >= 0 && randomNum < len(locations) {
		return locations[randomNum], nil
	}

	return models.Location{Address: "Texas"}, nil
}

type ApiBasedTool struct {
	wsclient *utils.WeatherStack
}

func GetApiBasedTool(wsclient *utils.WeatherStack) Tooler {
	return &ApiBasedTool{
		wsclient: wsclient,
	}
}

func (t *ApiBasedTool) GetLocation() (models.Location, error) {
	loc, err := utils.GetCurrentLocation()
	if err != nil {
		return models.Location{}, err
	}

	return models.Location{
		Address: loc.GenerateAddress(),
	}, nil
}

func (t *ApiBasedTool) GetCurrentWeather(loc models.Location) (models.Weather, error) {
	weather, err := t.wsclient.GetCurrentWeather(loc.ToString())
	if err != nil {
		return models.Weather{}, err
	}

	return *weather, nil
}
