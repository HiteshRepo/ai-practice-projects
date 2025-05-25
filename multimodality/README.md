# Multimodality

A Go CLI tool for generating, editing, and analyzing images using OpenAI's API. Supports image generation from text prompts, image inpainting (completion), and image vision (analysis) tasks.

## Features

- **Image Generation**: Create images from text prompts using OpenAI's image generation API.
- **Image Completion (Inpainting)**: Edit or complete images by providing an original and a masked image.
- **Image Vision**: Analyze and answer questions about images using OpenAI's vision models.

## Directory Structure

```
multimodality/
├── .env                # Environment variables (API keys, etc.)
├── go.mod, go.sum      # Go module files and dependencies
├── main.go             # Main CLI entry point
├── openai/             # OpenAI API client and image operation logic
└── test-files/         # Example images and test assets
```

## Setup

1. **Clone the repository** and navigate to the `multimodality` directory.

2. **Install dependencies** (requires Go 1.24+):

   ```bash
   go mod tidy
   ```

3. **Set up environment variables**:

   - Copy `.env` and set your OpenAI API key:
     ```
     export OPEN_API_KEY=your-openai-api-key
     ```

## Usage

Run the CLI with the desired action and flags:

### Image Generation

Generate an image from a text description:

```bash
go run main.go -action=image-gen -image-desc="An astronaut riding a bicycle on the moon"
```

### Image Completion (Inpainting)

Edit or complete an image using a mask:

```bash
go run main.go -action=image-complete -image-desc="Ancient Konark temple dedicated for Lord Surya (Sun) before it was destroyed." -image-path=./test-files/image.png -masked-image-path=./test-files/masked.png
```

### Image Vision (Analysis)

Ask a question about an image (local file or URL):

```bash
go run main.go -action=image-vision -query="What did the Ancient Konark temple dedicated for Lord Surya (Sun) look like before it was destroyed." -image-path=./test-files/image.png
```
or
```bash
go run main.go -action=image-vision -query="What did the Ancient Konark temple dedicated for Lord Surya (Sun) look like before it was destroyed." -image-url=https://ik.imagekit.io/1hhs6vx06v/konark.png
```

## Actions and Flags

| Flag                | Description                                                      |
|---------------------|------------------------------------------------------------------|
| `-action`           | One of: `image-gen`, `image-complete`, `image-vision`            |
| `-image-desc`       | Description for image generation or completion                   |
| `-image-path`       | Path to the input image file                                     |
| `-masked-image-path`| Path to the masked image file (for inpainting)                   |
| `-query`            | Question to ask about the image (for vision)                     |

## Dependencies

- [Go 1.24+](https://golang.org/)
- [OpenAI Go SDK](https://github.com/openai/openai-go)
- [caarlos0/env](https://github.com/caarlos0/env)
- [pkg/errors](https://github.com/pkg/errors)

## License

MIT License
