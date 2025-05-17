package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"pop-choice/constants"
	"pop-choice/models"
	openaipkg "pop-choice/openai"
	"pop-choice/supabase"

	"github.com/caarlos0/env"
	supa "github.com/nedpals/supabase-go"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"github.com/pkg/errors"
)

type envvars struct {
	OpenApiKey         string `env:"OPEN_API_KEY"`
	SupabaseApiKey     string `env:"SUPABASE_API_KEY"`
	SupabaseProjectUrl string `env:"SUPABASE_PROJECT_URL"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var envs envvars
	if err := env.Parse(&envs); err != nil {
		log.Fatalln("failed to parse env variables", errors.Wrap(err, "missing required env"))
	}

	actionFlag := flag.String(
		"action",
		"",
		"Allowed values: setup.")
	flag.Parse()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Printf("\nReceived signal: %s\n", sig)
		log.Println("Initiating graceful shutdown...")
		cancel()
	}()

	openaiClient := openaipkg.NewOpenAiClient(envs.OpenApiKey)
	supabaseClient := supabase.NewClient(envs.SupabaseProjectUrl, envs.SupabaseApiKey)
	scanner := bufio.NewScanner(os.Stdin)

	questionsTracker := 0
	multiUserQuestionsTracker := 0
	answersTracker := make([]string, len(constants.InitialListOfQuestions))
	allAnswersTracker := make([]string, 0)
	multiUserAnswersTracker := make([]string, len(constants.ExtraListOfQuestionsForMultiUser))
	reRun := false
	numOfUsersTracker := 1

	action := ""
	if actionFlag != nil {
		action = *actionFlag
	}

	switch action {
	case "setup":
		setup(ctx, openaiClient, supabaseClient)
	case "single-user":
		log.Println("Pop choice started. Press Ctrl+C to exit.")

		for {
			select {
			case <-ctx.Done():
				if err := scanner.Err(); err != nil {
					fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
				}

				log.Println("Pop choice ended. Hope you had a good time.")

				return
			default:
				if reRun {
					fmt.Print("> Do you wish to continue? yes/no: ")

					if reRunCheck(ctx, scanner) {
						questionsTracker = 0
					} else {
						return
					}
				}

				if questionsTracker < len(constants.InitialListOfQuestions) {
					fmt.Printf("> %s: ", constants.InitialListOfQuestions[questionsTracker])
					questionsTracker++
				}

				inputCh := make(chan string)
				errCh := make(chan error)
				go func() {
					if scanner.Scan() {
						inputCh <- scanner.Text()
					} else {
						if err := scanner.Err(); err != nil {
							errCh <- err
						} else {
							errCh <- fmt.Errorf("input closed")
						}
					}
				}()

				select {
				case <-ctx.Done():
					log.Println("Received shutdown signal during input. Exiting...")
					return
				case err := <-errCh:
					log.Fatalf("failed to read input: %v", err)
					return
				case input := <-inputCh:
					input = strings.TrimSpace(input)

					if len(input) == 0 {
						log.Println("You didn't enter anything. Try again!")
						continue
					}

					answersTracker[questionsTracker-1] = input

					if questionsTracker == len(constants.InitialListOfQuestions) {
						fmt.Println(generateResponse(
							ctx,
							openaiClient,
							supabaseClient,
							answersTracker,
							1,
							constants.PopChoiceSystemMessage))
						reRun = true
					}
				}
			}
		}
	case "multi-user":
		log.Println("Pop choice started. Press Ctrl+C to exit.")
		for {
			select {
			case <-ctx.Done():
				if err := scanner.Err(); err != nil {
					fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
				}

				log.Println("Pop choice ended. Hope you had a good time.")

				return
			default:
				if multiUserQuestionsTracker < len(constants.ExtraListOfQuestionsForMultiUser) {
					fmt.Printf("> %s: ", constants.ExtraListOfQuestionsForMultiUser[multiUserQuestionsTracker])
					multiUserQuestionsTracker++
				}

				if multiUserQuestionsTracker > len(constants.ExtraListOfQuestionsForMultiUser) &&
					questionsTracker < len(constants.InitialListOfQuestions) {
					fmt.Printf("> %s, %s: ",
						fmt.Sprintf("User-%d", numOfUsersTracker),
						constants.InitialListOfQuestions[questionsTracker])
					questionsTracker++
				}

				inputCh := make(chan string)
				errCh := make(chan error)
				go func() {
					if scanner.Scan() {
						inputCh <- scanner.Text()
					} else {
						if err := scanner.Err(); err != nil {
							errCh <- err
						} else {
							errCh <- fmt.Errorf("input closed")
						}
					}
				}()

				select {
				case <-ctx.Done():
					log.Println("Received shutdown signal during input. Exiting...")
					return
				case err := <-errCh:
					log.Fatalf("failed to read input: %v", err)
					return
				case input := <-inputCh:
					input = strings.TrimSpace(input)

					if len(input) == 0 {
						log.Println("You didn't enter anything. Try again!")
						continue
					}

					if multiUserQuestionsTracker > 0 && questionsTracker == 0 {
						multiUserAnswersTracker[multiUserQuestionsTracker-1] = input
					}

					if multiUserQuestionsTracker == len(constants.ExtraListOfQuestionsForMultiUser) {
						multiUserQuestionsTracker++
					}

					if questionsTracker > 0 {
						answersTracker[questionsTracker-1] = input
					}

					if multiUserQuestionsTracker >= len(constants.ExtraListOfQuestionsForMultiUser) &&
						questionsTracker == len(constants.InitialListOfQuestions) {
						numOfUsersTracker++

						if numOfUsers, err := strconv.Atoi(multiUserAnswersTracker[0]); err != nil && numOfUsers == numOfUsersTracker {
							fmt.Println(generateResponse(
								ctx,
								openaiClient,
								supabaseClient,
								allAnswersTracker,
								3,
								fmt.Sprintf(constants.MultiUserFilterTmpl, multiUserAnswersTracker[0], multiUserAnswersTracker[1])))
							return
						}

						allAnswersTracker = append(allAnswersTracker, []string{"\n", fmt.Sprintf("User-%d's interests", numOfUsersTracker), "\n"}...)
						allAnswersTracker = append(allAnswersTracker, answersTracker...)

						answersTracker = make([]string, len(constants.InitialListOfQuestions))
					}
				}
			}
		}
	}
}

func reRunCheck(
	ctx context.Context,
	scanner *bufio.Scanner,
) bool {
	inputCh := make(chan string)
	errCh := make(chan error)
	go func() {
		if scanner.Scan() {
			inputCh <- scanner.Text()
		} else {
			if err := scanner.Err(); err != nil {
				errCh <- err
			} else {
				errCh <- fmt.Errorf("input closed")
			}
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("Received shutdown signal during input. Exiting...")
		return false
	case err := <-errCh:
		log.Fatalf("failed to read input: %v", err)
		return false
	case input := <-inputCh:
		switch strings.ToLower(strings.TrimSpace(input)) {
		case "yes", "y":
			return true
		default:
			return false
		}
	}
}

func generateResponse(
	ctx context.Context,
	openaiClient openai.Client,
	supabaseClient *supa.Client,
	answersTracker []string,
	numberOfresponses int,
	systemPrompt string,
) []string {
	responses := make([]string, 0)

	systemMessage := openai.ChatCompletionMessageParamUnion{
		OfSystem: &openai.ChatCompletionSystemMessageParam{
			Content: openai.ChatCompletionSystemMessageParamContentUnion{
				OfString: openai.String(systemPrompt),
			},
		},
	}

	combinedInput := ""
	for _, a := range answersTracker {
		combinedInput += a + "\n"
	}

	res, err := openaiClient.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Model: "text-embedding-ada-002", // Default length of 1536 embeddings of array
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(combinedInput),
		},
	})
	if err != nil {
		log.Fatalln("failed to generate embeddings", err)
	}

	if res != nil && len(res.Data) > 0 {
		matchedPopChoiceMovies, err := supabase.InvokeMatchFunction(supabaseClient, constants.PopChoiceFunctionName, res.Data[0].Embedding, numberOfresponses)
		if err != nil {
			log.Fatalln("failed to match pop choice movies for query", err)
		}

		if len(matchedPopChoiceMovies) == 0 {
			log.Fatalln("invalid number of matching pop choice movies found")
		}

		for _, m := range matchedPopChoiceMovies {
			userMessage := openai.ChatCompletionMessageParamUnion{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String(
							fmt.Sprintf(
								"Context: %s, User interests: %s",
								m.Content, combinedInput),
						),
					},
				},
			}

			messages := []openai.ChatCompletionMessageParamUnion{systemMessage, userMessage}

			chatResp, err := openaiClient.Chat.Completions.New(
				ctx,
				openai.ChatCompletionNewParams{
					Messages:         messages,
					Model:            openai.ChatModelGPT4,
					Temperature:      param.NewOpt(constants.Temperature),
					PresencePenalty:  param.NewOpt(constants.PresencePenalty),
					FrequencyPenalty: param.NewOpt(constants.FrequencyPenalty),
				})
			if err != nil {
				log.Fatalln("failed to generate pop choice movies response", err)
			}

			responses = append(responses, chatResp.Choices[0].Message.Content)
		}
	}

	return responses
}

func setup(
	ctx context.Context,
	openaiClient openai.Client,
	supabaseClient *supa.Client,
) {
	log.Println("Starting setup.....")

	existingPopChoiceMovies := fetchExistingRows(supabaseClient, constants.PopChoiceTblName)

	log.Println("Fetched existing pop choice movies.....")

	for _, m := range constants.Movies {
		if _, ok := existingPopChoiceMovies[m.ToString()]; ok {
			log.Printf("Skipping insert for: %s\n", m.ToString())
			continue
		}

		movieVector := getEmbeddings(ctx, openaiClient, []string{m.ToString()})

		res, err := supabase.InsertDocument(constants.PopChoiceTblName, supabaseClient, movieVector[0])
		if err != nil {
			log.Fatalf("failed to insert embeddings for: '%s'\n: %v", movieVector[0].Content, err)
		}

		log.Printf("len of docs: %d\n", len(res))
	}
}

func fetchExistingRows(
	supabaseClient *supa.Client,
	tableName string) map[string]any {
	allDocs, err := supabase.ReadDocuments(tableName, supabaseClient)
	if err != nil {
		log.Fatalf("failed to fetch embeddings from: '%s'\n: %v", tableName, err)
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
