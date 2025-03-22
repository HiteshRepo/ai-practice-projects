package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/openai/openai-go"
	"github.com/pkg/errors"

	openaipkg "stock-price-predictor/openai"
	"stock-price-predictor/polygon"
)

const (
	SystemMessage = `You are a trading guru. Given data on share prices over past 3 days,
	write a report of no more that 150 words describing the stocks performance and recommend
	whether to buy or hold.`
)

type envvars struct {
	PolygonApiKey string `env:"POLYGON_API_KEY"`
	OpenApiKey    string `env:"OPEN_API_KEY"`
}

func main() {
	ctx := context.Background()

	var envs envvars
	if err := env.Parse(&envs); err != nil {
		log.Fatalln("failed to parse env variables", errors.Wrap(err, "missing required env"))
	}

	ticksFlag := flag.String("ticks", "MSFT", "comma separated stock names")
	flag.Parse()

	ticks := make([]string, 0)
	if ticksFlag != nil {
		for _, tick := range strings.Split(*ticksFlag, ",") {
			ticks = append(ticks, tick)
		}
	}

	openaiClient := openaipkg.NewOpenAiClient(envs.OpenApiKey)
	polygonClient := polygon.NewPolygonClient(envs.PolygonApiKey)

	totalTokenUsed := 0
	for _, tick := range ticks {
		dailyPrice, err := polygonClient.GetDailyPrices(tick, time.Now().Add(-72*time.Hour), time.Now())
		if err != nil {
			log.Fatalln(fmt.Sprintf("failed to fetch daily prices for %s: ", tick), err)
		}

		messages := []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(SystemMessage),
			openai.UserMessage(dailyPrice),
		}

		chatResp, err := openaiClient.Chat.Completions.New(
			ctx,
			openai.ChatCompletionNewParams{
				Messages: messages,
				Model:    openai.AudioModelGPT4oMiniTranscribe,
			})
		if err != nil {
			log.Fatalln(fmt.Sprintf("failed to generate stock report for %s: ", tick), err)
		}

		log.Printf("%s: %v\n\n", tick, chatResp.Choices[0].Message.Content)

		totalTokenUsed += int(chatResp.Usage.TotalTokens)
	}

	log.Printf("Total token used: %d\n", totalTokenUsed)
}
