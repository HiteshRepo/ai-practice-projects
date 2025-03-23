# Stock Price Predictor

A command-line tool that generates concise stock performance reports and buy/hold recommendations based on recent price data.

## Overview

Stock Price Predictor is a Go application that:
1. Fetches the last 3 days of stock price data from Polygon.io
2. Analyzes the data using OpenAI's GPT-4o Mini model
3. Supports two analysis approaches:
   - Zero-shot: Direct analysis with a system prompt
   - Few-shot: Analysis with example outputs to influence style
4. Generates a concise report (max 150 words) with a buy/hold recommendation

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

Run the application with default settings (analyzes Microsoft stock with zero-shot approach):

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

### Choosing Analysis Approach

The application supports two different approaches for generating reports:

```
# Zero-shot approach (default)
go run main.go -approach=zero

# Few-shot approach (uses examples to influence style)
go run main.go -approach=few
```

You can combine both flags:

```
go run main.go -ticks=NVDA,AMD -approach=few
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
3. The formatted data is sent to OpenAI's GPT-4o Mini model along with a prompt to analyze the stock's performance.
4. Depending on the approach selected:
   - **Zero-shot approach**: The data is sent directly to the model with a system message.
   - **Few-shot approach**: The data is sent along with example outputs to influence the style of the response.
5. The AI generates a concise report (under 150 words) that includes:
   - A summary of the stock's recent performance
   - A recommendation to buy or hold the stock
6. The report is displayed in the console.
7. The application tracks and reports the total number of tokens used for the OpenAI API calls (prompt tokens, completion tokens, and total).

## Example Output

### Zero-shot approach:
```
go run main.go -ticks=MSFT,TSLA -approach=zero

The stock MSFT opened at $385.74 on March 20th and closed at $386.84, with a high of $391.79 and a low of $383.28. The trading volume for the day was over 17 million shares, and the volume-weighted average price (VWAP) was $387.14. On March 21st, Microsoft opened at $383.21 and closed at $391.26, with a high of $391.74 and a low of $382.80. The trading volume was over 37 million shares.

Tesla (TSLA) opened at $233.34 on March 20th and closed at $236.26, with a high of $238.00 and a low of $230.05. The trading volume was over 98 million shares, and the VWAP was $234.31. On March 21st, Tesla opened at $234.99 and closed at $248.71, with a high of $249.52 and a low of $234.55. The trading volume was over 128 million shares.

From the data, Microsoft stocks have shown a gradual increase, while Tesla stocks have seen a significant rise. For investors looking for stability, Microsoft may be a good option. However, for those willing to take on higher risk, the increased activity in Tesla may present potential opportunities. It is recommended to hold Microsoft and to buy Tesla.

i/p token used: 431
o/p token used: 294
Total token used: 725
```

### Few-shot approach:
```
go run main.go -ticks=MSFT,TSLA -approach=few

MSFT STOCK PERFORMANCE:
Microsoft (MSFT) has had a volatile three-day period. The stock opened at $235.74 on day one, rose to $313.79 on day three, and closed at $313.96. On the last day, MSFT opened at $285.74, but experienced a sharp drop to $283.28, and closed at $283.96. The volume on day three was almost three times higher than the previous day. A volatile day for MSFT!

TSLA STOCK PERFORMANCE:
Tesla (TSLA) had a rollercoaster week. The stock opened at $223.34 on day one, but it closed at $245.46 due to a significant gain. However, the stock on day two declined by nearly $11, at $238.67, before a surge on the third day. On the last day, TSLA opened at $249.52, dipped to $244.55, and finally closed $248.71.

RECOMMENDATION:
Hold on to MSFT and TSLA for the long term. They have both been quite volatile, but the long-term outlook is promising.

i/p token used: 920
o/p token used: 237
Total token used: 1157
```

## Limitations

- The application only analyzes the past 3 days of stock data, which is a very short timeframe for investment decisions.
- The AI-generated recommendations should not be considered professional financial advice.
- API rate limits may apply depending on your Polygon.io and OpenAI subscription plans.
