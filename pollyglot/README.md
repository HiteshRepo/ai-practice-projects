# Pollyglot

Pollyglot is a command-line tool written in Go that leverages OpenAI's GPT-3.5-turbo model to translate any given content into a specified language. It is designed for developers and users who need quick, high-quality translations directly from the terminal.

## Features

- Translate any text to a target language using OpenAI's GPT-3.5-turbo.
- Simple CLI interface with flags for content and target language.
- Secure API key management via environment variables.
- Easily extensible and written in idiomatic Go.

## Requirements

- Go 1.24.1 or higher
- An OpenAI API key ([get one here](https://platform.openai.com/settings/organization/api-keys))

## Installation

1. **Clone the repository:**
   ```sh
   git clone <repo-url>
   cd pollyglot
   ```

2. **Install dependencies:**
   ```sh
   go mod tidy
   ```

3. **Build the binary:**
   ```sh
   go build -o pollyglot
   ```

## Environment Setup

Create a `.env` file in the `pollyglot` directory (or copy from `sample.env`) and add your OpenAI API key:

```
OPEN_API_KEY=your_openai_api_key_here
```

Alternatively, you can export the variable in your shell:

```sh
export OPEN_API_KEY=your_openai_api_key_here
```

## Usage

Run the CLI with the required flags:

```sh
./pollyglot -content="Hello, how are you?" -language="French"
```

- `-content`: The text you want to translate.
- `-language`: The target language for translation.

### Example

```sh
./pollyglot -content="Good morning, world!" -language="Spanish"
```

**Output:**
```
Buenos días, mundo!
```

## Project Structure

```
pollyglot/
├── go.mod
├── go.sum
├── main.go
├── openai/
│   └── client.go
├── sample.env
```

- `main.go`: CLI entry point and core logic.
- `openai/`: Contains the OpenAI client implementation.
- `sample.env`: Example environment variable file.

## Dependencies

- [openai-go](https://github.com/openai/openai-go)
- [caarlos0/env](https://github.com/caarlos0/env)
- [pkg/errors](https://github.com/pkg/errors)

## License

This project is provided as-is for educational and personal use.
