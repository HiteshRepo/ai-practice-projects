package main

import (
	"context"
	"flag"
	"log"
	"strings"

	"vector-embeddings/constants"
	"vector-embeddings/models"
	openaipkg "vector-embeddings/openai"
	"vector-embeddings/supabase"

	"github.com/caarlos0/env"
	supa "github.com/nedpals/supabase-go"
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

	actionFlag := flag.String("action", "insert", "Allowed values: insert, search. Default: search")
	searchQueryFlag := flag.String("query", "", "text for semantic search")
	flag.Parse()

	action := ""
	if actionFlag != nil {
		action = *actionFlag
	}

	query := ""
	if searchQueryFlag != nil {
		query = *searchQueryFlag
	}

	openaiClient := openaipkg.NewOpenAiClient(envs.OpenApiKey)
	supabaseClient := supabase.NewClient(envs.SupabaseProjectUrl, envs.SupabaseApiKey)

	switch action {
	case "insert":
		// go run main.go

		vectors := getEmbeddings(ctx, openaiClient, constants.Podcasts)
		allDocsMap := fetchAllDocumentsMap(supabaseClient)

		for _, v := range vectors {
			insertDocIfNotPresent(supabaseClient, v, allDocsMap)
		}
	case "search":
		// go run main.go -action=search -query="Jammin' in the Big Easy"
		// go run main.go -action=search -query="Decoding orca calls"
		// go run main.go -action=search -query="What can I listen to in half an hour?"

		if len(strings.TrimSpace(query)) == 0 {
			log.Fatalln("query cannot be empty for semantic search")
		}

		res, err := openaiClient.Embeddings.New(ctx, openai.EmbeddingNewParams{
			Model: "text-embedding-ada-002", // Default length of 1536 embeddings of array
			Input: openai.EmbeddingNewParamsInputUnion{
				OfString: openai.String(query),
			},
		})
		if err != nil {
			log.Fatalln("failed to generate embeddings", err)
		}

		if res != nil && len(res.Data) > 0 {
			matchedDocs, err := supabase.InvokeMatchDocumentsFunction(supabaseClient, res.Data[0].Embedding)
			if err != nil {
				log.Fatalln("failed to match documents for query", err)
			}

			if len(matchedDocs) == 0 {
				log.Fatalln("mo matching docs found")
			}

			for _, md := range matchedDocs {
				log.Printf("matched doc: %s, \nsimilarity score: %v\n", md.Content, md.Similarity)
			}
		} else {
			log.Fatalln("failed to generate embeddings for query", err)
		}
	}

}

func insertDocIfNotPresent(
	supabaseClient *supa.Client,
	v models.Vector,
	allDocsMap map[string]any,
) {
	res, err := supabase.ReadDocumentByContent(constants.DocumentsTblName, supabaseClient, v.Content)
	if err != nil {
		log.Fatalf("failed to find embeddings for: '%s'\n: %v", v.Content, err)
	}

	found := len(res) > 0

	if len(res) == 0 {
		_, found = allDocsMap[v.Content]
	}

	if !found {
		res, err := supabase.InsertDocument(constants.DocumentsTblName, supabaseClient, v)
		if err != nil {
			log.Fatalf("failed to insert embeddings for: '%s'\n: %v", v.Content, err)
		}

		log.Printf("len of docs: %d\n", len(res))
	}
}

func fetchAllDocumentsMap(supabaseClient *supa.Client) map[string]any {
	allDocs, err := supabase.ReadDocuments(constants.DocumentsTblName, supabaseClient)
	if err != nil {
		log.Fatalf("failed to find embeddings for: '%s'\n: %v", "random content", err)
	}

	allDocsMap := make(map[string]any)
	for _, d := range allDocs {
		allDocsMap[d.Content] = nil
	}

	return allDocsMap
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
