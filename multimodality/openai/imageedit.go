package openai

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/openai/openai-go"
)

func EditImage(
	ctx context.Context,
	openaiClient openai.Client,
	imageDesc string,
	maskedImagePath, actualImagePath string,
) {
	actualImage, err := fileReader(actualImagePath)
	if err != nil {
		log.Fatal(err)
	}

	maskedImage, err := fileReader(maskedImagePath)
	if err != nil {
		log.Fatal(err)
	}

	responseFormat := openai.ImageEditParamsResponseFormatB64JSON // default: url

	res, err := openaiClient.Images.Edit(
		ctx,
		openai.ImageEditParams{
			Prompt:         imageDesc,
			Model:          "dall-e-2", // at the moment image edit is available only in dall-e-2
			ResponseFormat: responseFormat,
			Image: openai.ImageEditParamsImageUnion{
				OfFile: actualImage,
			},
			Mask: maskedImage,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	switch responseFormat {
	case openai.ImageEditParamsResponseFormatB64JSON:
		b64JsonToPng(res.Data[0].B64JSON)
	}
}

func fileReader(
	filePath string) (*bytes.Reader, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	imageReader := bytes.NewReader(b)
	if imageReader == nil {
		return nil, fmt.Errorf("failed to create image reader")
	}

	return imageReader, nil
}
