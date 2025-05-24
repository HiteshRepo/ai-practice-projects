# reAct

reAct is a Go-based command-line tool that leverages OpenAI's API to generate personalized activity ideas based on your current location and weather. It supports multiple versions of activity suggestion logic (currently `v1`, `v2`, and `v3`) and is easily extensible.

## Features

- **AI-powered suggestions:** Uses OpenAI to generate creative activity ideas.
- **Weather and location aware:** Tailors suggestions to your current context.
- **Versioned logic:** Switch between different suggestion algorithms (`v1`, `v2`, `v3`). Default is `v3`.
- **Configurable via environment variables.**

## Version Overview

- **v1:** Basic activity suggestions using OpenAI, based on user query and weather.
- **v2:** Improved context handling and more personalized suggestions by refining prompt structure and response parsing.
- **v3:** Enhanced activity recommendations with advanced parsing, better error handling, and more robust logic for diverse scenarios.

## Directory Structure

```
reAct/
├── .env                # Environment variables (API keys, etc.)
├── constants/          # Project constants
├── go.mod, go.sum      # Go module files and dependencies
├── main.go             # Entry point for the CLI tool
├── models/             # Data models
├── openai/             # OpenAI API client logic
├── tools/              # Utility tools
├── utils/              # Helper utilities (e.g., parsers)
└── versions/           # Versioned activity suggestion logic
```

## Installation

1. **Clone the repository:**
   ```bash
   git clone <repo-url>
   cd reAct
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Set up your environment variables:**
   - Create a `.env` file or export your OpenAI API key in your shell:
     ```
     export OPEN_API_KEY=your-openai-api-key
     ```
   - **Do not commit your real API key to version control.**

## Usage

Run the CLI with your desired version and query:

```bash
go run main.go -version=v3 -query="Give me a list of activity ideas based on my current location and weather"
```

- `-version`: Selects the logic version (`v1`, `v2`, or `v3`). Default is `v3`.
- `-query`: Your prompt or question (e.g., about activities, weather, etc.).

### Examples

```bash
go run main.go -version=v1 -query="Give me a list of activity ideas based on my current location and weather"
```

```bash
go run main.go -version=v3 -query="Give me a list of activity ideas based on my current location and weather"
```

## Environment Variables

- `OPEN_API_KEY`: Your OpenAI API key (required).
- `WEATHER_STACK_API_KEY`: Your WeatherStack API key (required for weather-based suggestions).

You can set these in a `.env` file or export them in your shell.

## Dependencies

- [Go 1.24+](https://golang.org/)
- [OpenAI Go SDK](https://github.com/openai/openai-go)
- [caarlos0/env](https://github.com/caarlos0/env)
- [pkg/errors](https://github.com/pkg/errors)
