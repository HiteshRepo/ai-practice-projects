package main

import (
	"context"
	"log"

	"vector-embeddings/constants"
	"vector-embeddings/models"
	openaipkg "vector-embeddings/openai"
	"vector-embeddings/supabase"

	"github.com/caarlos0/env"
	"github.com/openai/openai-go"
	"github.com/pkg/errors"
)

type envvars struct {
	OpenApiKey         string `env:"OPEN_API_KEY"`
	SupabaseApiKey     string `env:"SUPABASE_API_KEY"`
	SupabaseProjectUrl string `env:"SUPABASE_PROJECT_URL"`
}

func main() {
	ctx := context.Background()

	var envs envvars
	if err := env.Parse(&envs); err != nil {
		log.Fatalln("failed to parse env variables", errors.Wrap(err, "missing required env"))
	}

	openaiClient := openaipkg.NewOpenAiClient(envs.OpenApiKey)
	supabaseClient := supabase.NewClient(envs.SupabaseProjectUrl, envs.SupabaseApiKey)

	vectors := getEmbeddings(ctx, openaiClient, constants.Podcasts)

	for _, v := range vectors {
		res, err := supabase.InsertDocument(constants.DocumentsTblName, supabaseClient, v)
		if err != nil {
			log.Fatalf("failed to insert embeddings for: '%s'\n: %v", v.Content, err)
		}

		log.Printf("len of docs: %d\n", len(res))
	}
}

func getEmbeddings(
	ctx context.Context,
	openaiClient openai.Client,
	sentences []string) []models.Vector {
	vectors := make([]models.Vector, 0)

	for _, s := range sentences {
		res, err := openaiClient.Embeddings.New(ctx, openai.EmbeddingNewParams{
			Model: "text-embedding-ada-002", // Default length of 1536 embeddings of array
			Input: openai.EmbeddingNewParamsInputUnion{
				OfString: openai.String(s),
			},
		})
		if err != nil {
			log.Fatalln("failed to generate embeddings", err)
		}

		vectors = append(vectors, models.Vector{
			Content:   s,
			Embedding: res.Data[0].Embedding,
		})
	}

	return vectors
}
