# Image Generator

A command-line tool that generates images from text descriptions using OpenAI's DALL-E models.

## Overview

Image Generator is a Go application that:
1. Takes a text description as input
2. Connects to OpenAI's API
3. Uses DALL-E 3 (default) to generate an image based on the description
4. Saves the generated image as a PNG file locally

## Prerequisites

- Go 1.24 or later
- API key for:
  - [OpenAI](https://platform.openai.com/) - for image generation

## Installation

1. Clone the repository:
   ```
   git clone <repository-url>
   cd image-generator
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

2. Edit the `.env` file and add your API key:
   ```
   OPEN_API_KEY=your_openai_api_key_here
   ```

   - Get your OpenAI API key from: https://platform.openai.com/settings/organization/api-keys

## Usage

### Basic Usage

Run the application with a text description of the image you want to generate:

```
go run main.go -image-prompt "A serene mountain landscape with a lake at sunset"
```

The generated image will be saved as `output.png` in the current directory.

### Building the Application

To build an executable:

```
go build -o image-generator
```

Then run it:

```
./image-generator -image-prompt "A futuristic city with flying cars"
```

## How It Works

1. The application takes a text description through the `-image-prompt` flag.
2. It connects to OpenAI's API using your provided API key.
3. The description is sent to OpenAI's DALL-E 3 model.
4. The model generates an image based on the description.
5. The application receives the image data in base64 JSON format.
6. The image is decoded and saved as `output.png` in the current directory.

## Features

- Uses DALL-E 3 by default (can be modified in the code to use DALL-E 2)
- Generates 1024x1024 pixel images
- Uses "vivid" style by default (can be changed to "natural" in the code)
- Saves images locally as PNG files

## Limitations

- The application currently saves all images with the same filename (`output.png`), overwriting previous generations.
- API rate limits and costs may apply depending on your OpenAI subscription plan.
- Image generation quality depends on the clarity and specificity of your text description.
