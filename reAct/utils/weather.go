package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"react/constants"
	"react/models"
	"time"
)

type WeatherStack struct {
	client *http.Client
	apiKey string
}

func GetWeatherStackClient(apiKey string) *WeatherStack {
	return &WeatherStack{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},

		apiKey: apiKey,
	}
}

func (ws *WeatherStack) GetCurrentWeather(loc string) (*models.Weather, error) {
	resp, err := ws.client.Get(
		fmt.Sprintf(constants.WeatherStackAPITmpl, ws.apiKey, loc))
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api request failed with status: %s", resp.Status)

	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	var weather models.WeatherResponse

	err = json.Unmarshal([]byte(body), &weather)
	if err != nil {
		return nil, fmt.Errorf("error deseriallizing weather response: %v", err)
	}

	return models.ConvertWeatherResponseToWeather(weather), nil
}
