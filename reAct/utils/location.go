package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"react/constants"
	"react/models"
	"strings"
	"time"
)

func getPublicIP(client *http.Client) (string, error) {
	for _, service := range constants.FetchIPAPIs {
		ip, err := fetchIP(client, service)
		if err == nil && ip != "" {
			return ip, nil
		}
	}

	return "", fmt.Errorf("failed to get IP from any service")
}

func fetchIP(client *http.Client, url string) (string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ip := strings.TrimSpace(string(body))
	return ip, nil
}

func GetCurrentLocation() (*models.LocationInfo, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	ip, err := getPublicIP(client)
	if err != nil {
		return nil, fmt.Errorf("error fetching ip: %v", err)

	}

	resp, err := client.Get(fmt.Sprintf("%s/%s", constants.IPToLocationAPI, ip))
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

	var li models.LocationInfo
	if err := json.Unmarshal(body, &li); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return &li, nil
}
