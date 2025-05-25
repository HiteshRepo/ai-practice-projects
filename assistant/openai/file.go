package openai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
)

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

func UploadFile(
	ctx context.Context,
	openaiClient openai.Client,
	filePath, fileName, purpose string,
) (string, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("file does not exist: %s", filePath)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	namedReader := &NamedReader{
		reader: bytes.NewReader(data),
		name:   fileName,
	}

	uploadedFile, err := openaiClient.Files.New(
		ctx,
		openai.FileNewParams{
			File: namedReader,
			// needs to be assistants to be used by assistant API
			Purpose: openai.FilePurpose(purpose),
		},
	)
	if err != nil {
		return "", err
	}

	return uploadedFile.ID, nil
}

func IsFileDataUploaded(
	ctx context.Context,
	client openai.Client,
	filePath, fileName, purpose string) (bool, string, error) {
	uploadedFiles, err := client.Files.List(ctx, openai.FileListParams{
		Purpose: param.NewOpt(purpose),
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
			Purpose: param.NewOpt(purpose),
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
