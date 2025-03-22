# Stock Price Predictor

A command-line tool that generates concise stock performance reports and buy/hold recommendations based on recent price data.

## Overview

Stock Price Predictor is a Go application that:
1. Fetches the last 3 days of stock price data from Polygon.io
2. Analyzes the data using OpenAI's GPT model
3. Generates a concise report (max 150 words) with a buy/hold recommendation

## Prerequisites

- Go 1.24 or later
- API keys for:
  - [Polygon.io](https://polygon.io/) - for stock price data
  - [OpenAI](https://platform.openai.com/) - for generating the analysis

## Installation

1. Clone the repository:
   ```
   git clone <repository-url>
   cd stock-price-predictor
   ```

2. Install dependencies:
   ```
   go mod download
   ```

## Configuration

1. Create a `.env` file in the project root directory based on the provided `sample.env`:
   ```
   cp sample.env .env
   ```

2. Edit the `.env` file and add your API keys:
   ```
   POLYGON_API_KEY=your_polygon_api_key_here
   OPEN_API_KEY=your_openai_api_key_here
   ```

   - Get your Polygon.io API key from: https://polygon.io/dashboard/keys
   - Get your OpenAI API key from: https://platform.openai.com/settings/organization/api-keys

## Usage

### Basic Usage

Run the application with default settings (analyzes Microsoft stock):

```
go run main.go
```

### Analyzing Specific Stocks

Specify one or more stock tickers using the `-ticks` flag:

```
# Analyze a single stock
go run main.go -ticks AAPL

# Analyze multiple stocks (comma-separated, no spaces)
go run main.go -ticks AAPL,GOOGL,AMZN,TSLA
```

### Building the Application

To build an executable:

```
go build -o stock-predictor
```

Then run it:

```
./stock-predictor -ticks NVDA,AMD
```

## How It Works

1. The application fetches daily price data (open, high, low, close, volume) for the specified stocks over the past 3 days from Polygon.io.
2. It formats this data into a readable table.
3. The formatted data is sent to OpenAI's GPT model along with a prompt to analyze the stock's performance.
4. The AI generates a concise report (under 150 words) that includes:
   - A summary of the stock's recent performance
   - A recommendation to buy or hold the stock
5. The report is displayed in the console.
6. The application tracks and reports the total number of tokens used for the OpenAI API calls.

## Example Output

```
AAPL: Apple's stock has shown moderate volatility over the past three days. 
Opening at $170.25, it experienced a slight dip to $168.90 before recovering to close at $171.15 on the final day. 
The trading volume has remained consistent, indicating stable investor interest. 
The stock has demonstrated resilience despite market fluctuations, with a positive upward trend forming by the end of the period. 
Given the stock's recovery pattern and the company's strong market position, I recommend a BUY for Apple shares at this time.

Total token used: 245
```

## Limitations

- The application only analyzes the past 3 days of stock data, which is a very short timeframe for investment decisions.
- The AI-generated recommendations should not be considered professional financial advice.
- API rate limits may apply depending on your Polygon.io and OpenAI subscription plans.
