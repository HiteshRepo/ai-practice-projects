# reAct

reAct is a Go-based command-line tool that leverages OpenAI's API to generate personalized activity ideas based on your current location and weather. It supports multiple versions of activity suggestion logic and is easily extensible.

## Features

- **AI-powered suggestions:** Uses OpenAI to generate creative activity ideas.
- **Weather and location aware:** Tailors suggestions to your current context.
- **Versioned logic:** Switch between different suggestion algorithms (`v1`, `v2`).
- **Configurable via environment variables.**

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
   - Copy `.env` and set your OpenAI API key:
     ```
     export OPEN_API_KEY=your-openai-api-key
     ```
   - You can also set this variable in your shell profile.

## Usage

Run the CLI with your desired version and query:

```bash
go run main.go -version=v2 -query="Give me a list of activity ideas based on my current location and weather"
```

- `-version`: Selects the logic version (`v1` or `v2`). Default is `v2`.
- `-query`: Your prompt or question (e.g., about activities, weather, etc.).

### Example

```bash
go run main.go -version=v1 -query="What can I do today in New York if it's raining?"
```

## Environment Variables

- `OPEN_API_KEY`: Your OpenAI API key (required).

You can set this in a `.env` file or export it in your shell.

## Dependencies

- [Go 1.24+](https://golang.org/)
- [OpenAI Go SDK](https://github.com/openai/openai-go)
- [caarlos0/env](https://github.com/caarlos0/env)
- [pkg/errors](https://github.com/pkg/errors)
