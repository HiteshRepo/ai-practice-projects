package polygon

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	PolygonBaseURL = "https://api.polygon.io"
)

type PolygonClient struct {
	APIKey     string
	HTTPClient *http.Client
}

func NewPolygonClient(apiKey string) *PolygonClient {
	return &PolygonClient{
		APIKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *PolygonClient) GetTickerDetails(symbol string) (*TickerResponse, error) {
	endpoint := fmt.Sprintf("%s/v3/reference/tickers/%s", PolygonBaseURL, url.PathEscape(symbol))

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	q := req.URL.Query()
	q.Add("apiKey", c.APIKey)
	req.URL.RawQuery = q.Encode()

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned error: %s - %s", resp.Status, string(body))
	}

	var tickerResp TickerResponse
	if err := json.NewDecoder(resp.Body).Decode(&tickerResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &tickerResp, nil
}

func (c *PolygonClient) GetDailyPrices(symbol string, from, to time.Time) (string, error) {
	endpoint := fmt.Sprintf("%s/v2/aggs/ticker/%s/range/1/day/%s/%s",
		PolygonBaseURL,
		url.PathEscape(symbol),
		from.Format("2006-01-02"),
		to.Format("2006-01-02"))

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	q := req.URL.Query()
	q.Add("apiKey", c.APIKey)
	q.Add("adjusted", "true")
	req.URL.RawQuery = q.Encode()

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API returned error: %s - %s", resp.Status, string(body))
	}

	var priceResp StockPriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&priceResp); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	return priceResp.ToString(), nil
}
