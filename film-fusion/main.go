package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/caarlos0/env"
	"github.com/openai/openai-go"
	"github.com/pkg/errors"

	"filmfusion/loader"
	openaipkg "filmfusion/openai"
)

const (
	MovieArtTextSystemMessage      = "You are a movie enthusiast with great writing skills. When provided with a movie name and art style, generate a single line exiting message using the context."
	MovieArtImageSystemMessageTmpl = `An imaginative poster inspired by the movie "%s", rendered in the "%s" art style.`
	UserQueryTmpl                  = "MovieName: %s, ArtStyle: %s"
)

var ArtStyles = []string{
	"art deco",
	"impressionism",
	"expressionism",
	"surrealism",
	"cubism",
	"cyberpunk",
	"abstract",
	"pop art",
	"minimalism",
	"futurism",
	"neoclassicism",
	"romanticism"}

func GetArtStyleOptions() string {
	artStyleOptions := ""

	for _, a := range ArtStyles {
		artStyleOptions += a + ", "
	}

	return artStyleOptions[0 : len(artStyleOptions)-2]
}

type envvars struct {
	OpenApiKey string `env:"OPEN_API_KEY"`
}

// go run main.go -film-name="Saving Private Ryan" -art-style=expressionism
// go run main.go -film-name="Inception" -art-style=cyberpunk
func main() {
	ctx := context.Background()

	var envs envvars
	if err := env.Parse(&envs); err != nil {
		log.Fatalln("failed to parse env variables", errors.Wrap(err, "missing required env"))
	}

	openaiClient := openaipkg.NewOpenAiClient(envs.OpenApiKey)

	filmNameFlag := flag.String("film-name", "", "name of a film")
	artStyleFlag := flag.String("art-style", "", fmt.Sprintf("valid options: %s", GetArtStyleOptions()))
	flag.Parse()

	filmName := ""
	if filmNameFlag != nil {
		filmName = *filmNameFlag
	}

	artStyle := ""
	if artStyleFlag != nil {
		artStyle = *artStyleFlag
	}

	if len(filmName) == 0 {
		flag.Usage()
		return
	}

	if len(artStyle) == 0 {
		flag.Usage()
		return
	}

	resp, err := openaiClient.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Model: "gpt-4",
			Messages: []openai.ChatCompletionMessageParamUnion{
				{
					OfSystem: &openai.ChatCompletionSystemMessageParam{
						Content: openai.ChatCompletionSystemMessageParamContentUnion{
							OfString: openai.String(MovieArtTextSystemMessage),
						},
					},
				},
				{
					OfUser: &openai.ChatCompletionUserMessageParam{
						Content: openai.ChatCompletionUserMessageParamContentUnion{
							OfString: openai.String(fmt.Sprintf(UserQueryTmpl, filmName, artStyle)),
						},
					},
				},
			},
		},
	)
	if err != nil {
		log.Fatalln(err)
	}

	content := resp.Choices[0].Message.Content

	loader := loader.NewLoader(content)
	loader.Start()

	imageResp, err := openaiClient.Images.Generate(
		ctx,
		openai.ImageGenerateParams{
			Prompt:         fmt.Sprintf(MovieArtImageSystemMessageTmpl, filmName, artStyle),
			Model:          "dall-e-3", // default: dall-e-2
			ResponseFormat: openai.ImageGenerateParamsResponseFormatURL,
			Size:           openai.ImageGenerateParamsSize1024x1792,
			Quality:        openai.ImageGenerateParamsQualityStandard,
			Style:          openai.ImageGenerateParamsStyleVivid,
		},
	)
	if err != nil {
		loader.Fail(fmt.Sprintf("Failed to generate image: %s", err.Error()))
	}

	loader.Complete(fmt.Sprintf("Check image at: %s", imageResp.Data[0].URL))
}
