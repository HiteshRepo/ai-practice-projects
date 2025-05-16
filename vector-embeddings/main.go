package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

	"vector-embeddings/constants"
	"vector-embeddings/langchain"
	"vector-embeddings/models"
	openaipkg "vector-embeddings/openai"
	"vector-embeddings/supabase"

	"github.com/caarlos0/env"
	supa "github.com/nedpals/supabase-go"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"github.com/pkg/errors"
)

const (
	PodcastsSystemMessage = `You are an enthusiastic podcast expert who loves recommending podcasts to people. You will be given two pieces of information - some context about podcasts episodes and a question. Your main job is to formulate a short answer to the question using the provided context. If you are unsure and cannot find the answer in the context, say, "Sorry, I don't know the answer." Please do not make up the answer.`
	UserMessageTmpl       = `Context: %s, Question: %s`

	MoviesSystemMessage        = `You are an enthusiastic movie expert who loves recommending movies to people. You will be given two pieces of information - some context about podcasts episodes and a question. Your main job is to formulate a short answer to the question using the provided context. If you are unsure and cannot find the answer in the context, say, "Sorry, I don't know the answer." Please do not make up the answer.`
	UserMessageMovieSearchTmpl = `Context: %s, Question: %s, NumberOfMovies: %d`

	Temperature      = 1.1
	PresencePenalty  = 0.0
	FrequencyPenalty = 0.0

	SpaceSplitter = " "
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

	actionFlag := flag.String(
		"action",
		"",
		"Allowed values: insert-docs, search-docs, search-n-chat-docs, chunk-n-insert-movies.")
	searchQueryFlag := flag.String(
		"query",
		"",
		"text for semantic search")
	matchesFlag := flag.Int(
		"matches",
		1,
		"number of matches to be used in your query")
	flag.Parse()

	action := ""
	if actionFlag != nil {
		action = *actionFlag
	}

	if len(action) == 0 {
		log.Fatal("mandatory flag `action` is not provided")
	}

	query := ""
	if searchQueryFlag != nil {
		query = *searchQueryFlag
	}

	matches := -1
	if matchesFlag != nil {
		matches = *matchesFlag
	}

	openaiClient := openaipkg.NewOpenAiClient(envs.OpenApiKey)
	supabaseClient := supabase.NewClient(envs.SupabaseProjectUrl, envs.SupabaseApiKey)

	switch action {
	case "insert-docs":
		// go run main.go -action=insert-docs

		allDocsMap := fetchExistingRows(supabaseClient, constants.DocumentsTblName)

		for _, p := range constants.Podcasts {
			if _, ok := allDocsMap[p]; ok {
				continue
			}

			docVector := getEmbeddings(ctx, openaiClient, []string{p})
			if len(docVector) == 0 {
				log.Printf("Failed to generate doc vector for: (%s)\n", p)
			}

			res, err := supabase.InsertDocument(constants.DocumentsTblName, supabaseClient, docVector[0])
			if err != nil {
				log.Fatalf("failed to insert embeddings for: '%s'\n: %v", docVector[0].Content, err)
			}

			log.Printf("len of docs: %d\n", len(res))
		}

	case "search-docs":
		// go run main.go -action=search-docs -query="Jammin' in the Big Easy"
		// go run main.go -action=search-docs -query="Decoding orca calls"
		// go run main.go -action=search-docs -query="What can I listen to in half an hour?"

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
			matchedDocs, err := supabase.InvokeMatchFunction(supabaseClient, constants.MatchDocumentsFunctionName, res.Data[0].Embedding, 2)
			if err != nil {
				log.Fatalln("failed to match documents for query", err)
			}

			if len(matchedDocs) == 0 {
				log.Fatalln("no matching docs found")
			}

			for _, md := range matchedDocs {
				log.Printf("matched doc: %s, \nsimilarity score: %v\n", md.Content, md.Similarity)
			}
		} else {
			log.Fatalln("failed to generate embeddings for query", err)
		}

	case "search-n-chat-docs":
		// go run main.go -action=search-n-chat-docs -query="Jammin' in the Big Easy"
		// go run main.go -action=search-n-chat-docs -query="Decoding orca calls"
		// go run main.go -action=ssearch-n-chat-docs -query="What can I listen to in half an hour?"
		// go run main.go -action=search-n-chat-docs -query="An episode Elon Musk would enjoy"

		if len(strings.TrimSpace(query)) == 0 {
			log.Fatalln("query cannot be empty for semantic search & chat")
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
			matchedDocs, err := supabase.InvokeMatchFunction(supabaseClient, constants.MatchDocumentsFunctionName, res.Data[0].Embedding, 1)
			if err != nil {
				log.Fatalln("failed to match documents for query", err)
			}

			if len(matchedDocs) != 1 {
				log.Fatalln("invalid number of matching docs found")
			}

			messages := []openai.ChatCompletionMessageParamUnion{
				{
					OfSystem: &openai.ChatCompletionSystemMessageParam{
						Content: openai.ChatCompletionSystemMessageParamContentUnion{
							OfString: openai.String(PodcastsSystemMessage),
						},
					},
				},
			}

			messages = append(messages, openai.ChatCompletionMessageParamUnion{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String(fmt.Sprintf(UserMessageTmpl, matchedDocs[0].Content, query)),
					},
				},
			})

			chatResp, err := openaiClient.Chat.Completions.New(
				ctx,
				openai.ChatCompletionNewParams{
					Messages:         messages,
					Model:            openai.ChatModelGPT4,
					Temperature:      param.NewOpt(Temperature),
					PresencePenalty:  param.NewOpt(PresencePenalty),
					FrequencyPenalty: param.NewOpt(FrequencyPenalty),
				})
			if err != nil {
				log.Fatalln("failed to generate podcast response", err)
			}

			log.Println(chatResp.Choices[0].Message.Content)
		} else {
			log.Fatalln("failed to generate embeddings for query", err)
		}

	case "chunk-n-insert-movies":
		// go run main.go -action=chunk-n-insert-movies

		log.Println("chunking movie details....")
		chunks, err := langchain.SplitDocuments(SpaceSplitter, constants.Movies)
		if err != nil {
			log.Fatalln("failed to split documents", err)
		}
		log.Println("chunking movie details finished....")

		log.Println("fetching existing movie embeddings....")
		existingMovies := fetchExistingRows(supabaseClient, constants.MoviesTblName)
		log.Println("fetching existing movie embeddings finished....")

		log.Println("processing movie chunks....")
		for _, ch := range chunks {
			if _, ok := existingMovies[ch]; ok {
				continue
			}

			movieVector := getEmbeddings(ctx, openaiClient, []string{ch})
			log.Printf("generating movie chunk (%s) embedding finished....\n", movieVector[0].Content)

			res, err := supabase.InsertDocument(constants.MoviesTblName, supabaseClient, movieVector[0])
			if err != nil {
				log.Fatalf("failed to insert embeddings for: '%s'\n: %v", movieVector[0].Content, err)
			}

			log.Printf("inserting (%s) embeddings finished\n...", movieVector[0].Content)
			log.Printf("len of docs: %d\n", len(res))
		}

	case "query-movie":
		// go run main.go -action=query-movie -query="Which movie can I take my child to?" -matches=3
		// go run main.go -action=query-movie -query="I feel like having a good laugh"
		// go run main.go -action=query-movie -query="Which movie will give me an adrenaline rush?" -matches=3
		// go run main.go -action=query-movie -query="What's the highest rated movie?"
		// go run main.go -action=query-movie -query="The movie with that actor from Castaway"

		if len(strings.TrimSpace(query)) == 0 {
			log.Fatalln("query cannot be empty for semantic search & chat")
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
			matchedDocs, err := supabase.InvokeMatchFunction(supabaseClient, constants.MatchMoviesFunctionName, res.Data[0].Embedding, matches)
			if err != nil {
				log.Fatalln("failed to match movies for query", err)
			}

			if len(matchedDocs) < matches {
				log.Fatalln("invalid number of matching movies found")
			}

			combinedMatchResult := ""
			for _, md := range matchedDocs {
				combinedMatchResult = fmt.Sprintf("%s\n%s", combinedMatchResult, md.Content)
			}

			messages := []openai.ChatCompletionMessageParamUnion{
				{
					OfSystem: &openai.ChatCompletionSystemMessageParam{
						Content: openai.ChatCompletionSystemMessageParamContentUnion{
							OfString: openai.String(MoviesSystemMessage),
						},
					},
				},
			}

			messages = append(messages, openai.ChatCompletionMessageParamUnion{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String(fmt.Sprintf(UserMessageMovieSearchTmpl, combinedMatchResult, query, matches)),
					},
				},
			})

			chatResp, err := openaiClient.Chat.Completions.New(
				ctx,
				openai.ChatCompletionNewParams{
					Messages:         messages,
					Model:            openai.ChatModelGPT4,
					Temperature:      param.NewOpt(Temperature),
					PresencePenalty:  param.NewOpt(PresencePenalty),
					FrequencyPenalty: param.NewOpt(FrequencyPenalty),
				})
			if err != nil {
				log.Fatalln("failed to generate movies response", err)
			}

			log.Println(chatResp.Choices[0].Message.Content)
		} else {
			log.Fatalln("failed to generate embeddings for query", err)
		}
	}

}

func fetchExistingRows(
	supabaseClient *supa.Client,
	tableName string) map[string]any {
	allDocs, err := supabase.ReadDocuments(tableName, supabaseClient)
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

func insertRowIfNotPresent(
	supabaseClient *supa.Client,
	tableName string,
	v models.Vector,
) {
	res, err := supabase.ReadDocumentByContent(tableName, supabaseClient, v.Content)
	if err != nil {
		log.Fatalf("failed to find embeddings for: '%s'\n: %v", v.Content, err)
	}

	found := len(res) > 0
	if !found {
		res, err := supabase.InsertDocument(tableName, supabaseClient, v)
		if err != nil {
			log.Fatalf("failed to insert embeddings for: '%s'\n: %v", v.Content, err)
		}

		log.Printf("len of docs: %d\n", len(res))
	}
}
