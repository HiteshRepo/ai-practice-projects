package tools

import "react/models"

type Tooler interface {
	GetCurrentWeather(models.Location) models.Weather
	GetLocation() models.Location
}

type HardCodedTool struct{}

func GetNewHardCodedTool() Tooler {
	return &HardCodedTool{}
}

func (t *HardCodedTool) GetCurrentWeather(_ models.Location) models.Weather {
	return models.Weather{
		Temperature: "35",
		Unit:        "C",
		Forecast:    "Sunny",
	}
}

func (t *HardCodedTool) GetLocation() models.Location {
	// an example real location for accurate result
	return models.Location{Address: "Delta Square, Bhubaneswar, Odisha, India"}
}
