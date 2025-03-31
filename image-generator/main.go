package main

import (
	"context"
	"encoding/base64"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/caarlos0/env"
	"github.com/openai/openai-go"
	"github.com/pkg/errors"

	openaipkg "image-generator/openai"
)

type ImageStyle string

const (
	DefaultImageStyle ImageStyle = "default"
	GhibliImageStyle  ImageStyle = "ghibli"
)

type envvars struct {
	OpenApiKey string `env:"OPEN_API_KEY"`
}

func main() {
	ctx := context.Background()

	var envs envvars
	if err := env.Parse(&envs); err != nil {
		log.Fatalln("failed to parse env variables", errors.Wrap(err, "missing required env"))
	}

	imagePromptFlag := flag.String("image-prompt", "", "describe the image you want to generate")
	imageStyleFlag := flag.String("image-style", "default", "Supported values: default, ghibli")
	imagePathFlag := flag.String("image-path", "", "Location of image")
	flag.Parse()

	if imagePromptFlag == nil {
		log.Fatalln("image-prompt flag is required")
	}

	imageStyle := DefaultImageStyle
	if imageStyleFlag != nil {
		imageStyle = ImageStyle(*imageStyleFlag)
	}

	imagePath := ""
	if imagePathFlag != nil {
		imagePath = *imagePathFlag
	}

	openaiClient := openaipkg.NewOpenAiClient(envs.OpenApiKey)

	var (
		imageResp *openaipkg.ImageGenerationResponse
		err       error
	)

	switch imageStyle {
	case GhibliImageStyle:
		respFormat := openai.ImageEditParamsResponseFormatB64JSON

		imageResp, err = openaipkg.GenerateGhibliStyleImage(ctx, openaiClient, imagePath, respFormat)
		if err != nil {
			log.Fatalln("failed to generate image", err)
		}

		switch respFormat {
		case openai.ImageEditParamsResponseFormatURL:
			log.Println("image generated successfully", imageResp.URL)
		case openai.ImageEditParamsResponseFormatB64JSON:
			b64JsonToPng(imageResp.B64Json)
		}
	default:
		imagePrompt := *imagePromptFlag
		respFormat := openai.ImageGenerateParamsResponseFormatB64JSON

		imageResp, err = openaipkg.GenerateImage(ctx, openaiClient, imagePrompt, respFormat)
		if err != nil {
			log.Fatalln("failed to generate image", err)
		}

		switch respFormat {
		case openai.ImageGenerateParamsResponseFormatURL:
			log.Println("image generated successfully", imageResp.URL)
		case openai.ImageGenerateParamsResponseFormatB64JSON:
			b64JsonToPng(imageResp.B64Json)
		}
	}

}

func b64JsonToPng(b64Json string) error {
	b64data := b64Json

	// If the base64 string contains metadata (like "data:image/png;base64,"), remove it
	if i := strings.Index(b64data, ","); i != -1 {
		b64data = b64data[i+1:]
	}

	imageData, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		return err
	}

	err = os.WriteFile("output.png", imageData, 0644)
	if err != nil {
		return err
	}

	return nil
}
