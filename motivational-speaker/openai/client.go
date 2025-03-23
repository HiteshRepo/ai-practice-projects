package openai

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
)

const (
	FineTunedModel = "gpt-3.5-turbo"
)

type UploadDataResponse struct {
	FileID string `json:"id"`
	Status string `json:"status"`
}

type FineTunedResponse struct {
	UploadDataResponse UploadDataResponse
}

func NewOpenAiClient(openApiKey string) openai.Client {
	return openai.NewClient(
		option.WithAPIKey(openApiKey),
	)
}

func IsTrainingDataUploaded(
	ctx context.Context,
	client openai.Client,
	trainingDataFilePath string,
	purpose openai.FilePurpose) (bool, string, error) {
	fileName := filepath.Base(trainingDataFilePath)

	uploadedFiles, err := client.Files.List(ctx, openai.FileListParams{
		Purpose: param.NewOpt(string(purpose)),
	})
	if err != nil {
		return false, "", err
	}

	for _, file := range uploadedFiles.Data {
		if file.Filename == fileName {
			return true, file.ID, nil
		}
	}

	after := ""
	if len(uploadedFiles.Data) > 0 {
		after = uploadedFiles.Data[len(uploadedFiles.Data)-1].ID
	}

	for {
		if !uploadedFiles.HasMore {
			break
		}

		nextPage, err := client.Files.List(ctx, openai.FileListParams{
			Purpose: param.NewOpt(string(purpose)),
			After:   param.NewOpt(after),
		})
		if err != nil {
			return false, "", err
		}

		if len(nextPage.Data) == 0 {
			break
		}

		for _, file := range nextPage.Data {
			if file.Filename == fileName {
				return true, file.ID, nil
			}
		}

		after = nextPage.Data[len(nextPage.Data)-1].ID
		uploadedFiles = nextPage
	}

	return false, "", nil
}

// This helps with providing a name to the file being uploaded to openai.
// With just io.Reader, the file is uploaded with name 'anonymous_file'.
type NamedReader struct {
	reader io.Reader
	name   string
}

func (nr *NamedReader) Read(p []byte) (n int, err error) {
	return nr.reader.Read(p)
}

func (nr *NamedReader) Name() string {
	return nr.name
}

func UploadTrainingData(
	ctx context.Context,
	client openai.Client,
	trainingDataFilePath string,
	purpose openai.FilePurpose) (*FineTunedResponse, error) {
	fileName := filepath.Base(trainingDataFilePath)

	data, err := os.ReadFile(trainingDataFilePath)
	if err != nil {
		return nil, err
	}

	namedReader := &NamedReader{
		reader: bytes.NewReader(data),
		name:   fileName,
	}

	uploadedFile, err := client.Files.New(ctx, openai.FileNewParams{
		File:    namedReader,
		Purpose: purpose,
	})
	if err != nil {
		return nil, err
	}

	return &FineTunedResponse{
		UploadDataResponse: UploadDataResponse{
			FileID: uploadedFile.ID,
			Status: string(uploadedFile.Status),
		},
	}, nil
}

func StartAndAwaitFineTuneJob(
	ctx context.Context,
	client openai.Client,
	fileID string) (*openai.FineTuningJob, error) {
	var fineTuneJob *openai.FineTuningJob = nil

	resp, err := client.FineTuning.Jobs.New(ctx, openai.FineTuningJobNewParams{
		Model:        FineTunedModel,
		TrainingFile: fileID,
	})
	if err != nil {
		return nil, err
	}

	for {
		if resp.Status == openai.FineTuningJobStatusSucceeded {
			fineTuneJob = resp
			break
		}

		time.Sleep(5 * time.Second)

		resp, err = client.FineTuning.Jobs.Get(ctx, resp.ID)
		if err != nil {
			return nil, err
		}

		log.Printf("Job(%s) status: %s\n", resp.ID, resp.Status)
	}

	if fineTuneJob != nil {
		return fineTuneJob, nil
	}

	return nil, errors.New("fine tune job failed")
}

func RetrieveLatestFineTuneJobByTrainingFileID(
	ctx context.Context,
	client openai.Client,
	trainingFileID string) (*openai.FineTuningJob, error) {
	var latestJob *openai.FineTuningJob = nil

	jobs, err := client.FineTuning.Jobs.List(ctx, openai.FineTuningJobListParams{})
	if err != nil {
		return nil, err
	}

	for _, job := range jobs.Data {
		if job.TrainingFile == trainingFileID {
			currJob, err := client.FineTuning.Jobs.Get(ctx, job.ID)
			if err != nil {
				return nil, err
			}

			if latestJob != nil &&
				currJob.CreatedAt > latestJob.CreatedAt {
				latestJob = currJob
			}

			if latestJob == nil {
				latestJob = currJob
			}
		}
	}

	after := ""
	if len(jobs.Data) > 0 {
		after = jobs.Data[len(jobs.Data)-1].ID
	}

	for {
		if !jobs.HasMore {
			break
		}

		nextPage, err := client.FineTuning.Jobs.List(ctx, openai.FineTuningJobListParams{
			After: param.NewOpt(after),
		})
		if err != nil {
			return nil, err
		}

		if len(nextPage.Data) == 0 {
			break
		}

		for _, job := range nextPage.Data {
			if job.TrainingFile == trainingFileID {
				currJob, err := client.FineTuning.Jobs.Get(ctx, job.ID)
				if err != nil {
					return nil, err
				}

				if latestJob != nil &&
					currJob.CreatedAt > latestJob.CreatedAt {
					latestJob = currJob
				}

				if latestJob == nil {
					latestJob = currJob
				}
			}
		}

		after = nextPage.Data[len(nextPage.Data)-1].ID
		jobs = nextPage
	}

	return latestJob, nil
}

func RetrieveFineTuneJobByID(
	ctx context.Context,
	client openai.Client,
	jobID string) (*openai.FineTuningJob, error) {
	return client.FineTuning.Jobs.Get(ctx, jobID)
}
