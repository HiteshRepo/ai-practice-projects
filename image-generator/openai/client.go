package openai

import (
	"bytes"
	"context"
	"fmt"
	"image-generator/imageops"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
)

const (
	DallE2 = "dall-e-2"
	DallE3 = "dall-e-3"
)

const (
	GhibliStyleImageDescription = "Convert the provided image into Ghibli style artwork"
)

type ImageGenerationResponse struct {
	// Does not get generated in the response for dall-e-2
	RevisedPrompt string `json:"revised_prompt"`
	// URL is valid for only 1 hour
	URL     string `json:"url"`
	B64Json string `json:"b64_json"`
}

func NewOpenAiClient(openApiKey string) openai.Client {
	return openai.NewClient(
		option.WithAPIKey(openApiKey),
	)
}

func GenerateImage(
	ctx context.Context,
	client openai.Client,
	imageDescription string,
	respFormat openai.ImageGenerateParamsResponseFormat) (*ImageGenerationResponse, error) {
	imgResp, err := client.Images.Generate(ctx, openai.ImageGenerateParams{
		Prompt:         imageDescription,
		Model:          DallE3,                 // defaults to dall-e-2
		N:              param.NewOpt(int64(1)), // defaults to 1
		Size:           openai.ImageGenerateParamsSize1024x1024,
		Style:          openai.ImageGenerateParamsStyleVivid, // defaults to vivid (other option: Natural)
		ResponseFormat: respFormat,                           // defaults to URL (other option: b64_json)
	})
	if err != nil {
		return nil, err
	}

	return &ImageGenerationResponse{
		URL:           imgResp.Data[0].URL,
		RevisedPrompt: imgResp.Data[0].RevisedPrompt,
		B64Json:       imgResp.Data[0].B64JSON,
	}, nil
}

func GenerateGhibliStyleImage(
	ctx context.Context,
	client openai.Client,
	imagePath string,
	respFormat openai.ImageEditParamsResponseFormat) (*ImageGenerationResponse, error) {
	b, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, err
	}

	imageReader := bytes.NewReader(b)
	if imageReader == nil {
		return nil, fmt.Errorf("failed to create image reader")
	}

	// Create the mask image (example: make a 100x100 area starting at (50,50) editable)
	maskBytes, err := imageops.CreateMaskImage(imagePath, 50, 50, 100, 100)
	if err != nil {
		return nil, err
	}

	maskReader := bytes.NewReader(maskBytes)
	if maskReader == nil {
		return nil, fmt.Errorf("failed to create mask image reader")
	}

	imgResp, err := client.Images.Edit(ctx, openai.ImageEditParams{
		Prompt:         GhibliStyleImageDescription,
		Image:          imageReader,
		Mask:           maskReader,
		Model:          DallE3,                 // defaults to dall-e-2
		N:              param.NewOpt(int64(1)), // defaults to 1
		Size:           openai.ImageEditParamsSize1024x1024,
		ResponseFormat: respFormat, // defaults to URL (other option: b64_json)
	})
	if err != nil {
		return nil, err
	}

	return &ImageGenerationResponse{
		URL:           imgResp.Data[0].URL,
		RevisedPrompt: imgResp.Data[0].RevisedPrompt,
		B64Json:       imgResp.Data[0].B64JSON,
	}, nil
}
