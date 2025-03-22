package polygon

import (
	"fmt"
	"strings"
	"time"
)

type TickerResponse struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
	Result struct {
		Ticker           string    `json:"ticker"`
		Name             string    `json:"name"`
		Market           string    `json:"market"`
		Locale           string    `json:"locale"`
		PrimaryExchange  string    `json:"primary_exchange"`
		Type             string    `json:"type"`
		Active           bool      `json:"active"`
		CurrencyName     string    `json:"currency_name"`
		CurrencySymbol   string    `json:"currency_symbol"`
		Description      string    `json:"description"`
		HomepageURL      string    `json:"homepage_url"`
		ListDate         string    `json:"list_date"`
		MarketCap        int64     `json:"market_cap"`
		PhoneNumber      string    `json:"phone_number"`
		ShareClassShares int64     `json:"share_class_shares_outstanding"`
		WeightedShares   int64     `json:"weighted_shares_outstanding"`
		SicCode          string    `json:"sic_code"`
		SicDescription   string    `json:"sic_description"`
		TickerRoot       string    `json:"ticker_root"`
		TickerSuffix     string    `json:"ticker_suffix"`
		LastUpdatedUTC   time.Time `json:"last_updated_utc"`
	} `json:"results"`
}

type StockPriceResponse struct {
	Status       string `json:"status"`
	RequestID    string `json:"request_id"`
	Ticker       string `json:"ticker"`
	QueryCount   int    `json:"queryCount"`
	ResultsCount int    `json:"resultsCount"`
	Adjusted     bool   `json:"adjusted"`
	Results      []struct {
		Volume         float64 `json:"v"`
		VolumeWeighted float64 `json:"vw"`
		Open           float64 `json:"o"`
		Close          float64 `json:"c"`
		High           float64 `json:"h"`
		Low            float64 `json:"l"`
		Timestamp      int64   `json:"t"`
		NumberOfTrades int     `json:"n"`
	} `json:"results"`
}

func (spr *StockPriceResponse) ToString() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Stock Price Data for %s\n", spr.Ticker))
	builder.WriteString(fmt.Sprintf("Status: %s\n", spr.Status))
	builder.WriteString(fmt.Sprintf("Request ID: %s\n", spr.RequestID))
	builder.WriteString(fmt.Sprintf("Results: %d of %d\n", spr.ResultsCount, spr.QueryCount))
	builder.WriteString(fmt.Sprintf("Adjusted: %t\n\n", spr.Adjusted))

	builder.WriteString("Date       | Open     | High     | Low      | Close    | Volume      | VWAP     | Trades\n")
	builder.WriteString("-----------|----------|----------|----------|----------|-------------|----------|--------\n")

	for _, result := range spr.Results {
		date := time.Unix(result.Timestamp/1000, 0).Format("2006-01-02")
		builder.WriteString(fmt.Sprintf(
			"%-10s | $%-7.2f | $%-7.2f | $%-7.2f | $%-7.2f | $%-7.2f | $%-7.2f | %d\n",
			date,
			result.Open,
			result.High,
			result.Low,
			result.Close,
			result.Volume,
			result.VolumeWeighted,
			result.NumberOfTrades,
		))
	}

	return builder.String()
}
