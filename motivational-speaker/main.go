package main

import (
	"context"
	"log"
	"time"

	"github.com/caarlos0/env"
	"github.com/openai/openai-go"
	"github.com/pkg/errors"

	openaipkg "motivational-speaker/openai"
)

const (
	TrainingDataFilePath = "./finetunedata.jsonl"
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

	openaiClient := openaipkg.NewOpenAiClient(envs.OpenApiKey)

	isFileAlreadyUploaded, fileID, err := openaipkg.IsTrainingDataUploaded(
		ctx,
		openaiClient,
		TrainingDataFilePath,
		openai.FilePurposeFineTune)
	if err != nil {
		log.Fatalln("failed to check if training data is uploaded", err)
	}

	if !isFileAlreadyUploaded {
		uploadDataResponse, err := openaipkg.UploadTrainingData(
			ctx,
			openaiClient,
			TrainingDataFilePath,
			openai.FilePurposeFineTune)
		if err != nil {
			log.Fatalln("failed to upload training data", err)
		}

		log.Println("uploaded training data", uploadDataResponse.UploadDataResponse.FileID)
		fileID = uploadDataResponse.UploadDataResponse.FileID
	} else {
		log.Println("training data already uploaded", fileID)
	}

	newFineTuneJobRequired := false
	fineTuneJob, err := openaipkg.RetrieveLatestFineTuneJobByTrainingFileID(
		ctx,
		openaiClient,
		fileID)
	if err != nil {
		log.Println("failed to retrieve fine tune job by training file name", err)
		newFineTuneJobRequired = true
	} else {
		log.Println("fine tune job retrieved", fineTuneJob.ID)

		if fineTuneJob.Status == openai.FineTuningJobStatusFailed {
			newFineTuneJobRequired = true
		}

		attempts := 0
		for {
			if fineTuneJob.Status == openai.FineTuningJobStatusSucceeded {
				log.Println("fine tune job in succeed state")
				break
			}

			fineTuneJob, err = openaipkg.RetrieveFineTuneJobByID(ctx, openaiClient, fineTuneJob.ID)
			if err != nil {
				log.Println("failed to retrieve fine tune job by id", err)

				newFineTuneJobRequired = true
				break
			}

			if fineTuneJob.Status == openai.FineTuningJobStatusFailed {
				newFineTuneJobRequired = true
				break
			}

			time.Sleep(10 * time.Second)
			log.Printf("Job(%s) status: %s\n", fineTuneJob.ID, fineTuneJob.Status)

			if attempts > 5 {
				log.Println("exhausted attempts to check fine tune job status, starting new job")

				newFineTuneJobRequired = true
				break
			}

			attempts += 1
		}
	}

	if newFineTuneJobRequired {
		log.Println("starting new fine tune job")

		fineTuneJob, err = openaipkg.StartAndAwaitFineTuneJob(ctx, openaiClient, fileID)
		if err != nil {
			log.Fatalln("failed to start/complete fine tune job", err)
		}

		log.Println("fine tune job completed")
	}

	log.Println("going to use fine tuned model", fineTuneJob.Model)

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage("I don't know what is the purpose of my life"),
	}

	chatResp, err := openaiClient.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Messages: messages,
			Model:    fineTuneJob.Model,
		})
	if err != nil {
		log.Fatalln("failed to generate motivation response", err)
	}

	log.Println(chatResp.Choices[0].Message.Content)
}
