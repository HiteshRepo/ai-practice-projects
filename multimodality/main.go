package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/caarlos0/env"
	"github.com/pkg/errors"

	openaipkg "multimodality/openai"
)

const (
	// Image details: Setting, Objects, Colors, Mood, Anything you envision
	// Lighting, Time of day, Angle, Distance

	BasicPrompt          = "An astronaut riding a bicycle on the moon"
	DetailedPrompt       = "A colorful image of an astronaut cycling on the moon, with a vibrant Earth in the background. Include glowing tire tracks, colorful alien plants, and crystals on the lunar surface"
	NonRevisedPromptTmpl = "I NEED to test how the tool works with extremely simple prompts. DO NOT add any detail, just use it AS-IS: %s"
)

type envvars struct {
	OpenApiKey string `env:"OPEN_API_KEY"`
}

// go run main.go -action=image-gen -image-desc="An astronaut riding a bicycle on the moon"
// go run main.go -action=image-complete -image-desc="Ancient Konark temple dedicated for Lord Surya (Sun) before it was destroyed." -image-path=./test-files/image.png -masked-image-path=./test-files/masked.png
// go run main.go -action=image-vision -query="What did the Ancient Konark temple dedicated for Lord Surya (Sun) look like before it was destroyed." -image-url=https://ik.imagekit.io/1hhs6vx06v/konark.png
// go run main.go -action=image-vision -query="What did the Ancient Konark temple dedicated for Lord Surya (Sun) look like before it was destroyed." -image-path=./test-files/image.png
func main() {
	ctx := context.Background()

	var envs envvars
	if err := env.Parse(&envs); err != nil {
		log.Fatalln("failed to parse env variables", errors.Wrap(err, "missing required env"))
	}

	actionFlag := flag.String("action", "", "supported values: image-gen, image-complete")
	imageDescFlag := flag.String("image-desc", "", "description of image you would like to generate")
	actualImageFlag := flag.String("image-path", "", "path of actual image")
	maskedImageFlag := flag.String("masked-image-path", "", "path of masked image")
	queryFlag := flag.String("query", "", "what is your question about the image")
	flag.Parse()

	action := ""
	if actionFlag != nil {
		action = *actionFlag
	}

	imageDesc := ""
	if imageDescFlag != nil {
		imageDesc = *imageDescFlag
	}

	query := ""
	if queryFlag != nil {
		query = *queryFlag
	}

	actualImagePath := ""
	if actualImageFlag != nil {
		actualImagePath = *actualImageFlag
	}

	maskedImagePath := ""
	if maskedImageFlag != nil {
		maskedImagePath = *maskedImageFlag
	}

	validate(action, imageDesc, maskedImagePath, actualImagePath, query)

	openaiClient := openaipkg.NewOpenAiClient(envs.OpenApiKey)

	switch action {
	case "image-gen":
		openaipkg.GenerateImage(ctx, openaiClient, imageDesc)

	case "image-complete":
		openaipkg.EditImage(ctx, openaiClient, imageDesc, maskedImagePath, actualImagePath)

	case "image-vision":
		openaipkg.VisualizeImage(ctx, openaiClient, query, actualImagePath)
	}

}

func validate(
	action, imageDesc, maskedImagePath, actualImagePath, query string) {
	if len(action) == 0 {
		flag.Usage()
		return
	}

	if action == "image-gen" && len(imageDesc) == 0 {
		flag.Usage()
		return
	}

	if action == "image-complete" &&
		(len(maskedImagePath) == 0 || len(actualImagePath) == 0) {
		flag.Usage()
		return
	}

	if action == "image-complete" {
		if _, err := os.Stat(actualImagePath); os.IsNotExist(err) {
			log.Fatalf("file does not exist: %s", actualImagePath)
		}

		if _, err := os.Stat(maskedImagePath); os.IsNotExist(err) {
			log.Fatalf("file does not exist: %s", maskedImagePath)
		}
	}

	if action == "image-vision" && (len(query) == 0 || len(actualImagePath) == 0) {
		flag.Usage()
		return
	}
}
