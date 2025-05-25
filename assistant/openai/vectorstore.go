package openai

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
)

func CreateVectorStore(
	ctx context.Context,
	openaiClient openai.Client,
	name string,
	fileIDs []string,
) (string, error) {
	vectorStore, err := openaiClient.VectorStores.New(
		ctx,
		openai.VectorStoreNewParams{
			Name:    openai.String(name),
			FileIDs: fileIDs,
		},
	)
	if err != nil {
		return "", err
	}

	return vectorStore.ID, nil
}

func IsVectorStoreCreated(
	ctx context.Context,
	client openai.Client,
	name string) (bool, string, error) {
	vectorStores, err := client.VectorStores.List(
		ctx,
		openai.VectorStoreListParams{})
	if err != nil {
		return false, "", err
	}

	for _, vs := range vectorStores.Data {
		if vs.Name == name {
			return true, vs.ID, nil
		}
	}

	after := ""
	if len(vectorStores.Data) > 0 {
		after = vectorStores.Data[len(vectorStores.Data)-1].ID
	}

	for {
		if !vectorStores.HasMore {
			break
		}

		nextPage, err := client.VectorStores.List(
			ctx,
			openai.VectorStoreListParams{
				After: param.NewOpt(after),
			})
		if err != nil {
			return false, "", err
		}

		if len(nextPage.Data) == 0 {
			break
		}

		for _, vs := range nextPage.Data {
			if vs.Name == name {
				return true, vs.ID, nil
			}
		}

		after = nextPage.Data[len(nextPage.Data)-1].ID
		vectorStores = nextPage
	}

	return false, "", nil
}
