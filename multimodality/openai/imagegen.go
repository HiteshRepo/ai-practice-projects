package openai

import (
	"context"
	"encoding/base64"
	"log"
	"os"
	"strings"

	"github.com/openai/openai-go"
)

func GenerateImage(
	ctx context.Context,
	openaiClient openai.Client,
	imageDesc string,
) {
	responseFormat := openai.ImageGenerateParamsResponseFormatB64JSON // default: url

	res, err := openaiClient.Images.Generate(
		ctx,
		openai.ImageGenerateParams{
			// N:      openai.Int(3), // Not supported in dall-e-3
			Prompt:         imageDesc,
			Model:          "dall-e-3", // default: dall-e-2
			ResponseFormat: responseFormat,
			Size:           openai.ImageGenerateParamsSize1024x1792,
			Quality:        openai.ImageGenerateParamsQualityHD,  // HD is double the cost of normal
			Style:          openai.ImageGenerateParamsStyleVivid, // Vivid (default): Hyper-real, Dramatic
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	switch responseFormat {
	case openai.ImageGenerateParamsResponseFormatB64JSON:
		b64JsonToPng(res.Data[0].B64JSON)
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
