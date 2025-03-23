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
	"github.com/openai/openai-go/packages/param"
	"github.com/pkg/errors"

	openaipkg "stock-price-predictor/openai"
	"stock-price-predictor/polygon"
)

const (
	SystemMessage = `You are a trading guru. Given data on share prices over past 3 days,
	write a report of no more that 150 words describing the stocks performance and recommend
	whether to buy or hold. 
	Please use examples, if provided, between ### to set the style of response while generating your report.`

	OutputExample1 = `OK baby, hold on tight! You are going to haate this! 
	Over the past three days, Tesla (TSLA) shares have plummetted. 
	The stock opened at $223.98 and closed at $202.11 on the third day, with some jumping around in the meantime. 
	This is a great time to buy, baby! But not a great time to sell! 
	But I'm not done! Apple (AAPL) stocks have gone stratospheric! This is a seriously hot stock right now. 
	They opened at $166.38 and closed at $182.89 on day three. 
	So all in all, I would hold on to Tesla shares tight if you already have them - they might bounce right back up and head to the stars! They are volatile stock, so expect the unexpected. 
	For APPL stock, how much do you need the money? Sell now and take the profits or hang on and wait for more! 
	If it were me, I would hang on because this stock is on fire right now!!! 
	Apple are throwing a Wall Street party and y'all invited!
	`

	OutputExample2 = `Apple (AAPL) is the supernova in the stock sky - it shot up from $150.22 to a jaw-dropping $175.36 by the close of day three. 
	We're talking about a stock that's hotter than a pepper sprout in a chilli cook-off, and it's showing no signs of cooling down! 
	If you're sitting on AAPL stock, you might as well be sitting on the throne of Midas. 
	Hold on to it, ride that rocket, and watch the fireworks, because this baby is just getting warmed up! 
	Then there's Meta (META), the heartthrob with a penchant for drama. It winked at us with an opening of $142.50, but by the end of the thrill ride, it was at $135.90, leaving us a little lovesick. 
	It's the wild horse of the stock corral, bucking and kicking, ready for a comeback. 
	META is not for the weak-kneed So, sugar, what's it going to be? For AAPL, my advice is to stay on that gravy train. 
	As for META, keep your spurs on and be ready for the rally.
	`
)

const (
	// Determines creativity/daring-ness of the response.
	Temperature = 1.1

	// Higher the value, more topics get covered in the response.
	// Lower the value, more focused the response.
	PresencePenalty = 0.0

	// Higher the value, less repetitive phrases used in the response.
	// Lower the value, more repetitive phrases 'maybe' used in the response.
	FrequencyPenalty = 0.0

	// Use carefully, as it can stop the response generation abruptly.
	// finish_reason = "length" as against "stop".
	// MaxTokens   = 150
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
	approachFlag := flag.String("approach", "zero", "expected values: 'zero', 'few'")
	flag.Parse()

	ticks := make([]string, 0)
	if ticksFlag != nil {
		for _, tick := range strings.Split(*ticksFlag, ",") {
			ticks = append(ticks, tick)
		}
	}

	approach := ""
	if approachFlag != nil {
		approach = *approachFlag
	}

	openaiClient := openaipkg.NewOpenAiClient(envs.OpenApiKey)
	polygonClient := polygon.NewPolygonClient(envs.PolygonApiKey)

	stockPriceDetails := ""
	for _, tick := range ticks {
		dailyPrice, err := polygonClient.GetDailyPrices(tick, time.Now().Add(-72*time.Hour), time.Now())
		if err != nil {
			log.Fatalln(fmt.Sprintf("failed to fetch daily prices for %s: ", tick), err)
		}

		stockPriceDetails += dailyPrice + "\n\n"
	}

	switch approach {
	case "zero":
		zeroShot(ctx, openaiClient, stockPriceDetails)
	case "few":
		fewShot(ctx, openaiClient, stockPriceDetails)
	}
}

func zeroShot(
	ctx context.Context,
	openaiClient openai.Client,
	stockPriceDetails string) {
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(SystemMessage),
		openai.UserMessage(stockPriceDetails),
	}

	generateReport(ctx, openaiClient, messages)
}

func fewShot(ctx context.Context,
	openaiClient openai.Client,
	stockPriceDetails string) {

	withExample1 := fmt.Sprintf("%s\n\n###\n%s###\n\n", stockPriceDetails, OutputExample1)
	withExample2 := fmt.Sprintf("%s\n\n###\n%s###\n\n", withExample1, OutputExample2)

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(SystemMessage),
		openai.UserMessage(withExample2),
	}

	generateReport(ctx, openaiClient, messages)
}

func generateReport(
	ctx context.Context,
	openaiClient openai.Client,
	messages []openai.ChatCompletionMessageParamUnion) {
	chatResp, err := openaiClient.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Messages:         messages,
			Model:            openai.AudioModelGPT4oMiniTranscribe,
			Temperature:      param.NewOpt(Temperature),
			PresencePenalty:  param.NewOpt(PresencePenalty),
			FrequencyPenalty: param.NewOpt(FrequencyPenalty),
		})
	if err != nil {
		log.Fatalln("failed to generate stock report", err)
	}

	log.Println(chatResp.Choices[0].Message.Content)
	log.Println()

	log.Printf("i/p token used: %d\n", int(chatResp.Usage.PromptTokens))
	log.Printf("o/p token used: %d\n", int(chatResp.Usage.CompletionTokens))
	log.Printf("Total token used: %d\n", int(chatResp.Usage.TotalTokens))
}
