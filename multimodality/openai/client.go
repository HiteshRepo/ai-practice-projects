package openai

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func NewOpenAiClient(openApiKey string) openai.Client {
	return openai.NewClient(
		option.WithAPIKey(openApiKey),
	)
}
