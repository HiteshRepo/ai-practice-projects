package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/caarlos0/env"
	"github.com/openai/openai-go"
	"github.com/pkg/errors"

	openaipkg "pollyglot/openai"
)

const (
	SystemMessage = `You are a multiple language expert.
	For a provided content and language, please do the translation.`

	UserMessageTmpl = `Translate the following content to %s: %s`

	GPT3_5ModelName = "gpt-3.5-turbo"
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

	contentFlag := flag.String("content", "", "pass in relevant content to translate to asked language")
	languageFlag := flag.String("language", "", "language to translate to")
	flag.Parse()

	content := ""
	if contentFlag != nil {
		content = *contentFlag
	}

	language := ""
	if languageFlag != nil {
		language = *languageFlag
	}

	openaiClient := openaipkg.NewOpenAiClient(envs.OpenApiKey)

	userMessage := fmt.Sprintf(UserMessageTmpl, language, content)
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(SystemMessage),
		openai.UserMessage(userMessage),
	}

	chatResp, err := openaiClient.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Messages: messages,
			Model:    GPT3_5ModelName,
		})
	if err != nil {
		log.Fatalln("failed to translate the message", err)
	}

	log.Println(chatResp.Choices[0].Message.Content)
}
