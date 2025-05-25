package main

import (
	"context"
	"flag"
	"log"

	"github.com/caarlos0/env"
	"github.com/openai/openai-go"
	"github.com/pkg/errors"

	"assistant/constants"
	openaipkg "assistant/openai"
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

	queryFlag := flag.String("query", "", "ask a question on movie recommendation")
	flag.Parse()

	query := ""
	if queryFlag != nil {
		query = *queryFlag
	}

	if len(query) == 0 {
		flag.Usage()
		return
	}

	openaiClient := openaipkg.NewOpenAiClient(envs.OpenApiKey)

	uploaded, fileID, err := openaipkg.IsFileDataUploaded(
		ctx,
		openaiClient,
		constants.MovieDetailsFilePath,
		constants.MovieDetailsFileName,
		constants.MovieAssistantPurpose)
	if err != nil {
		log.Fatal(err)
	}

	if !uploaded {
		log.Printf("uploading %s file...\n", constants.MovieDetailsFileName)

		fileID, err = openaipkg.UploadFile(
			ctx,
			openaiClient,
			constants.MovieDetailsFilePath,
			constants.MovieDetailsFileName,
			constants.MovieAssistantPurpose)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("FileID: %s\n", fileID)

	created, vsID, err := openaipkg.IsVectorStoreCreated(
		ctx,
		openaiClient,
		constants.MoviesVectorStoreName)
	if err != nil {
		log.Fatal(err)
	}

	if !created {
		log.Printf("creating %s vector store...\n", constants.MoviesVectorStoreName)

		vsID, err = openaipkg.CreateVectorStore(
			ctx,
			openaiClient,
			constants.MoviesVectorStoreName,
			[]string{fileID})
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("VectorStoreID: %s\n", vsID)

	assistant := movieAssistant(ctx, openaiClient, vsID)

	log.Println("Movie assistant ID: ", assistant.ID)

	thread, err := openaipkg.CreateThread(ctx, openaiClient, map[string]string{})
	if err != nil {
		log.Fatal(err)
	}
	defer openaipkg.DeleteThread(ctx, openaiClient, thread.ID)

	err = openaipkg.AddMessageToThread(
		ctx,
		openaiClient,
		thread.ID,
		query,
		openai.BetaThreadMessageNewParamsRoleUser)
	if err != nil {
		log.Fatal(err)
	}

	err = openaipkg.RunThread(ctx, openaiClient, thread.ID, assistant.ID, constants.MovieRunInstructions)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := openaipkg.GetThreadMessages(ctx, openaiClient, thread.ID)
	if err != nil {
		log.Fatal(err)
	}

	if len(msgs) == 0 {
		log.Fatal(errors.New("no messages could be retrieved"))
	}

	log.Println(msgs[0])
}

func movieAssistant(
	ctx context.Context,
	openaiClient openai.Client,
	vectorStoreID string,
) *openai.Assistant {
	created, assistant, err := openaipkg.IsAssistantCreated(
		ctx,
		openaiClient,
		constants.MovieAssistantName)
	if err != nil {
		log.Fatal(err)
	}

	if !created {
		assistant, err = openaipkg.CreateAssistant(
			ctx,
			openaiClient,
			constants.MovieAssistantInstructions,
			constants.MovieAssistantName,
			nil,
			&openai.BetaAssistantNewParamsToolResources{
				FileSearch: openai.BetaAssistantNewParamsToolResourcesFileSearch{
					VectorStoreIDs: []string{vectorStoreID},
				},
			},
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	return assistant
}

func testAssistant(
	ctx context.Context,
	openaiClient openai.Client,
) string {
	testAssistantName := "Math Tutor"

	created, assistant, err := openaipkg.IsAssistantCreated(
		ctx,
		openaiClient,
		testAssistantName)
	if err != nil {
		log.Fatal(err)
	}

	if !created {
		assistant, err = openaipkg.CreateAssistant(
			ctx,
			openaiClient,
			"You are a personal math tutor. When asked a question, write and run Python code to answer the question.",
			testAssistantName,
			[]openai.AssistantToolUnionParam{
				{
					OfCodeInterpreter: &openai.CodeInterpreterToolParam{},
				},
			},
			nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	return assistant.ID
}
