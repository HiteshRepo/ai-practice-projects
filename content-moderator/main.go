package main

import (
	"context"
	"flag"
	"log"

	"github.com/caarlos0/env"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"github.com/pkg/errors"

	"content-moderator/fileops"
	openaipkg "content-moderator/openai"
	"content-moderator/scraper"
)

type envvars struct {
	OpenApiKey string `env:"OPEN_API_KEY"`
}

type ContentTypes string

const (
	SingleLineText ContentTypes = "single_line_text"
	MultiLineText  ContentTypes = "multi_line_text"
	TextFromURL    ContentTypes = "text_from_url"
	TextFromFile   ContentTypes = "text_from_file"
)

func main() {
	ctx := context.Background()

	var envs envvars
	if err := env.Parse(&envs); err != nil {
		log.Fatalln("failed to parse env variables", errors.Wrap(err, "missing required env"))
	}

	contentFlag := flag.String("content", "", "pass in relevant content to check if it is safe or not")
	typeFlag := flag.String("type", "single_line_text", "Allowed values: single_line_text, multi_line_text, text_from_url, text_from_file. Default: single_line_text")
	urlLinkFlag := flag.String("url", "", "pass in the url to check if the content is safe or not")
	fileLocFlag := flag.String("file", "", "pass in the file location to check if the content is safe or not")
	flag.Parse()

	content := ""
	if contentFlag != nil {
		content = *contentFlag
	}

	typeOfContent := SingleLineText
	if typeFlag != nil {
		switch *typeFlag {
		case string(SingleLineText):
			typeOfContent = SingleLineText
		case string(MultiLineText):
			typeOfContent = MultiLineText
		case string(TextFromURL):
			typeOfContent = TextFromURL
		case string(TextFromFile):
			typeOfContent = TextFromFile
		default:
			log.Fatalln("invalid content type provided")
		}
	}

	urlLink := ""
	if urlLinkFlag != nil {
		urlLink = *urlLinkFlag
	}

	fileLoc := ""
	if fileLocFlag != nil {
		fileLoc = *fileLocFlag
	}

	openaiClient := openaipkg.NewOpenAiClient(envs.OpenApiKey)

	// TODO: Add support for multiple results
	resultLookupIndex := -1
	var inputContent openai.ModerationNewParamsInputUnion

	switch typeOfContent {
	case SingleLineText:
		resultLookupIndex = 0

		inputContent = openai.ModerationNewParamsInputUnion{
			OfString: param.NewOpt(content),
		}

	case TextFromURL:
		content, err := scraper.ScrapeURL(urlLink)
		if err != nil {
			log.Fatalln("failed to scrape URL:", err)
		}

		resultLookupIndex = 0
		inputContent = openai.ModerationNewParamsInputUnion{
			OfString: param.NewOpt(content),
		}

	case TextFromFile:
		content, err := fileops.ReadFile(fileLoc)
		if err != nil {
			log.Fatalln("failed to scrape URL:", err)
		}

		resultLookupIndex = 0
		inputContent = openai.ModerationNewParamsInputUnion{
			OfString: param.NewOpt(content),
		}
	}

	res, err := openaiClient.Moderations.New(ctx, openai.ModerationNewParams{
		Input: inputContent,
	})
	if err != nil {
		log.Fatalln("failed to check content moderation:", err)
	}

	if res.Results[resultLookupIndex].Flagged {
		cats := getFlaggedCategories(res, resultLookupIndex)
		log.Println("Your content is flagged as unsafe in categories:", cats)
	} else {
		log.Println("Your content is safe")
	}
}

func getFlaggedCategories(
	res *openai.ModerationNewResponse,
	resultLookupIndex int,
) []string {
	cats := make([]string, 0)

	if res.Results[resultLookupIndex].Categories.Harassment {
		cats = append(cats, "Harassment")
	}

	if res.Results[resultLookupIndex].Categories.HarassmentThreatening {
		cats = append(cats, "Harassment Threatening")
	}

	if res.Results[resultLookupIndex].Categories.Hate {
		cats = append(cats, "Hate")
	}

	if res.Results[resultLookupIndex].Categories.HateThreatening {
		cats = append(cats, "Hate Threatening")
	}

	if res.Results[resultLookupIndex].Categories.SelfHarm {
		cats = append(cats, "Self Harm")
	}

	if res.Results[resultLookupIndex].Categories.SelfHarmInstructions {
		cats = append(cats, "Self Harm Instructions")
	}

	if res.Results[resultLookupIndex].Categories.SelfHarmIntent {
		cats = append(cats, "Self Harm Intent")
	}

	if res.Results[resultLookupIndex].Categories.Sexual {
		cats = append(cats, "Sexual")
	}

	if res.Results[resultLookupIndex].Categories.SexualMinors {
		cats = append(cats, "Sexual Minors")
	}

	if res.Results[resultLookupIndex].Categories.Violence {
		cats = append(cats, "Violence")
	}

	if res.Results[resultLookupIndex].Categories.ViolenceGraphic {
		cats = append(cats, "Violence Graphic")
	}

	if res.Results[resultLookupIndex].Categories.Illicit {
		cats = append(cats, "Illicit")
	}

	if res.Results[resultLookupIndex].Categories.IllicitViolent {
		cats = append(cats, "Illicit Violent")
	}

	return cats
}
