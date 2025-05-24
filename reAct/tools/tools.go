package tools

import (
	"encoding/json"
	"log"
	"react/models"
	"react/utils"

	"github.com/openai/openai-go"
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
	GetCurrentWeather(models.Location) (models.Weather, error)
	GetLocation() (models.Location, error)
}

func InvokeValidAction(
	wsClient *utils.WeatherStack,
	action models.Action) string {
	switch action.FunctionName {
	case "getCurrentWeather":
		log.Println("calling function getCurrentWeather")

		if len(action.Arguments) != 1 {
			log.Fatalln("invalid number of arguments for getCurrentWeather action")
		}

		weather, err := GetApiBasedTool(wsClient).GetCurrentWeather(models.Location{Address: action.Arguments[0]})
		if err != nil {
			log.Fatalln("failed to get weather: ", err)
		}

		return weather.ToString()

	case "getLocation":
		log.Println("calling function getLocation")

		if len(action.Arguments) != 1 {
			log.Fatalln("invalid number of arguments for getLocation action")
		}

		if !utils.IsEmpty(action.Arguments[0]) {
			log.Fatalln("getLocation action requires no args")
		}

		loc, err := GetApiBasedTool(wsClient).GetLocation()
		if err != nil {
			log.Fatalln("failed to get location: ", err)
		}

		return loc.ToString()
	}

	log.Fatalf("unknown action invoked: %s", action.FunctionName)

	return ""
}

func ActionsFromResponseToolCalls(
	toolCalls []openai.ChatCompletionMessageToolCall) ([]models.Action, error) {
	actions := make([]models.Action, 0)

	for _, tool := range toolCalls {
		action := models.Action{
			FunctionName: tool.Function.Name,
			Arguments:    make([]string, 0),
			ToolCallID:   tool.ID,
		}

		switch tool.Function.Name {
		case "getCurrentWeather":
			weatherArg := models.WeatherInputs{}
			err := json.Unmarshal([]byte(tool.Function.Arguments), &weatherArg)
			if err != nil {
				return nil, err
			}

			action.Arguments = append(action.Arguments, weatherArg.Address)

		case "getLocation":
			action.Arguments = append(action.Arguments, "NONE")
		}

		actions = append(actions, action)
	}

	return actions, nil
}
