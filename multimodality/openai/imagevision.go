package openai

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/openai/openai-go"
)

func VisualizeImage(
	ctx context.Context,
	openaiClient openai.Client,
	query string,
	imagePath string,
) {
	encodedImage, err := base64ImageFile(imagePath)
	if err != nil {
		log.Fatal(err)
	}

	image := fmt.Sprintf("data:image/png;base64,%s", encodedImage)

	res, err := openaiClient.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Model: "gpt-4o",
			Messages: []openai.ChatCompletionMessageParamUnion{
				{
					OfUser: &openai.ChatCompletionUserMessageParam{
						Content: openai.ChatCompletionUserMessageParamContentUnion{
							OfArrayOfContentParts: []openai.ChatCompletionContentPartUnionParam{
								{
									OfText: &openai.ChatCompletionContentPartTextParam{
										Text: query,
									},
								},
								{
									OfImageURL: &openai.ChatCompletionContentPartImageParam{
										ImageURL: openai.ChatCompletionContentPartImageImageURLParam{
											URL: image,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res.Choices[0].Message.Content)
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
