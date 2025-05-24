# Film Fusion

Film Fusion is a command-line tool that generates AI-powered art inspired by your favorite movies and a chosen art style. It leverages OpenAI's GPT-4 for creative prompt generation and DALL-E 3 for image synthesis, producing unique artwork based on your input.

## Features

- Generate a creative, single-line prompt based on a movie name and art style.
- Use OpenAI's GPT-4 to craft imaginative prompts.
- Create vivid images with DALL-E 3 based on the generated prompt.
- Supports a variety of art styles (see below).
- Simple command-line interface.

## Requirements

- Go 1.24+
- An OpenAI API key with access to GPT-4 and DALL-E 3

## Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/hiteshrepo/film-fusion.git
   cd film-fusion
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Configure your OpenAI API key:**
   - Copy `.env` and set your OpenAI API key:
     ```
     export OPEN_API_KEY=sk-...
     ```
   - You can set this in your shell or use a tool like [direnv](https://direnv.net/).

## Usage

Run the tool with the required flags:

```bash
go run main.go -film-name="Inception" -art-style=cyberpunk
```

- `-film-name` (required): The name of the movie.
- `-art-style` (required): The desired art style (see options below).

If either flag is missing, usage instructions will be displayed.

### Art Style Options

You can choose from the following art styles:

- art deco
- impressionism
- expressionism
- surrealism
- cubism
- cyberpunk
- abstract
- pop art
- minimalism
- futurism
- neoclassicism
- romanticism

### Example

```bash
go run main.go -film-name="Saving Private Ryan" -art-style=expressionism
```

The tool will:
1. Generate a creative prompt based on your input.
2. Display a loading animation.
3. Generate an image using DALL-E 3.
4. Output a URL to view the generated image.

## Project Structure

- `main.go` - Entry point and CLI logic.
- `loader/` - Loading animation and progress display.
- `openai/` - OpenAI API client logic.
- `utils/` - Utility functions (e.g., image handling).

## Acknowledgements

- [OpenAI GPT-4](https://platform.openai.com/docs/models/gpt-4)
- [OpenAI DALL-E 3](https://platform.openai.com/docs/guides/images/usage)
- [caarlos0/env](https://github.com/caarlos0/env) for environment variable parsing
- [pkg/errors](https://github.com/pkg/errors) for error handling
