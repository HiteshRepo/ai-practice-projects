package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/caarlos0/env"
	"github.com/openai/openai-go"
	"github.com/pkg/errors"

	openaipkg "billpricecalculator/openai"
)

const (
	Cheese1ImagePath = "images/cheese-1.jpeg"
	Cheese2ImagePath = "images/cheese-2.jpeg"
	MenuImagePath    = "images/menu.png"
)

var UseCases = map[string][]string{
	"What's the difference between these two types of cheese?":                                      {Cheese1ImagePath, Cheese2ImagePath},
	"I want to order one of each item on this menu for my company party. How much would that cost?": {MenuImagePath},
}

type envvars struct {
	OpenApiKey string `env:"OPEN_API_KEY"`
}

func main() {
	ctx := context.Background()

	var envs envvars
	if err := env.Parse(&envs); err != nil {
		log.Fatalln("failed to parse env variables", errors.Wrap(err, "missing required env"))
	}

	openaiClient := openaipkg.NewOpenAiClient(envs.OpenApiKey)

	for uc, images := range UseCases {
		imageContents := make([]openai.ChatCompletionContentPartUnionParam, 0)
		for _, img := range images {
			extension := filepath.Ext(img)

			encodedImage, err := base64ImageFile(img)
			if err != nil {
				log.Fatal(err)
			}

			image := fmt.Sprintf("data:image/%s;base64,%s", extension, encodedImage)

			imageContents = append(
				imageContents,
				openai.ChatCompletionContentPartUnionParam{
					OfImageURL: &openai.ChatCompletionContentPartImageParam{
						ImageURL: openai.ChatCompletionContentPartImageImageURLParam{
							URL: image,
						},
					},
				},
			)
		}

		resp := VisualizeImage(ctx, openaiClient, uc, imageContents)
		log.Printf("Usecase: %s, Response: %s\n\n\n", uc, resp)
	}
}

func VisualizeImage(
	ctx context.Context,
	openaiClient openai.Client,
	query string,
	imageContents []openai.ChatCompletionContentPartUnionParam,
) string {
	contents := []openai.ChatCompletionContentPartUnionParam{
		{
			OfText: &openai.ChatCompletionContentPartTextParam{
				Text: query,
			},
		},
	}

	contents = append(contents, imageContents...)

	res, err := openaiClient.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Model: "gpt-4o",
			Messages: []openai.ChatCompletionMessageParamUnion{
				{
					OfUser: &openai.ChatCompletionUserMessageParam{
						Content: openai.ChatCompletionUserMessageParamContentUnion{
							OfArrayOfContentParts: contents,
						},
					},
				},
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	return res.Choices[0].Message.Content
}

func base64ImageFile(
	filePath string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	b, err := os.ReadFile(filepath.Join(cwd, filePath))
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(b)
	return encoded, nil
}
