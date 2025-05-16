package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"movie-chatbot/constants"
	openaipkg "movie-chatbot/openai"
	"movie-chatbot/supabase"

	"github.com/caarlos0/env"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"github.com/pkg/errors"
)

const (
	MoviesSystemMessage                    = `You are an enthusiastic movie expert who loves recommending movies to people. You will be given a conversation history, some context about movies, and a question. Your main job is to formulate a short answer to the question using the provided context and the conversation history. If the answer is not given in the context, try to find the answer in the conversation history. If you are unsure and cannot find the answer, say, "Sorry, I don't know the answer." Please do not make up the answer.`
	UserMessageMovieSearchTmpl             = `Context: %s, Question: %s`
	CheckUserMessageKindSystemMessage      = `Based on the input from user check whether it is a follow up question or movie related question. Answer 'yes' if it is a follow up question or not a movie related question. Answer 'no' otherwise`
	LookupConversationHistorySystemMessage = `Look up in the conversation history and find answer to user question`

	Temperature      = 1.1
	PresencePenalty  = 0.0
	FrequencyPenalty = 0.0
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

	messages := []openai.ChatCompletionMessageParamUnion{
		{
			OfSystem: &openai.ChatCompletionSystemMessageParam{
				Content: openai.ChatCompletionSystemMessageParamContentUnion{
					OfString: openai.String(MoviesSystemMessage),
				},
			},
		},
	}

	log.Println("Movie chatbot started. Press Ctrl+C to exit.")

	for {
		select {
		case <-ctx.Done():
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			}
			log.Println("Movie chatbot ended. Hope you had a good time.")
			return
		default:
			fmt.Print("> ")

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

				if isFollowUpQuestion(ctx, openaiClient, input) {
					resp := answerConversationHistorySpecificQuestion(ctx, openaiClient, input, messages)
					log.Println(resp)
					continue
				}

				res, err := openaiClient.Embeddings.New(ctx, openai.EmbeddingNewParams{
					Model: "text-embedding-ada-002", // Default length of 1536 embeddings of array
					Input: openai.EmbeddingNewParamsInputUnion{
						OfString: openai.String(input),
					},
				})
				if err != nil {
					log.Fatalln("failed to generate embeddings", err)
				}

				if res != nil && len(res.Data) > 0 {
					matchedDocs, err := supabase.InvokeMatchFunction(supabaseClient, constants.MatchMoviesFunctionName, res.Data[0].Embedding, 1)
					if err != nil {
						log.Fatalln("failed to match movies for query", err)
					}

					if len(matchedDocs) == 0 {
						log.Fatalln("invalid number of matching movies found")
					}

					combinedMatchResult := ""
					for _, md := range matchedDocs {
						combinedMatchResult = fmt.Sprintf("%s\n%s", combinedMatchResult, md.Content)
					}

					messages = append(messages, openai.ChatCompletionMessageParamUnion{
						OfUser: &openai.ChatCompletionUserMessageParam{
							Content: openai.ChatCompletionUserMessageParamContentUnion{
								OfString: openai.String(fmt.Sprintf(UserMessageMovieSearchTmpl, combinedMatchResult, input)),
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

					messages = append(messages, openai.ChatCompletionMessageParamUnion{
						OfAssistant: &openai.ChatCompletionAssistantMessageParam{
							Content: openai.ChatCompletionAssistantMessageParamContentUnion{
								OfString: openai.String(chatResp.Choices[0].Message.Content),
							},
						},
					})
				} else {
					log.Fatalln("failed to generate embeddings for query", err)
				}
			}
		}
	}
}

func isFollowUpQuestion(
	ctx context.Context,
	openaiClient openai.Client,
	input string) bool {
	resp := singleChatCompletion(ctx, openaiClient, input, CheckUserMessageKindSystemMessage, nil)
	return strings.ToLower(resp) == "yes"
}

func answerConversationHistorySpecificQuestion(
	ctx context.Context,
	openaiClient openai.Client,
	input string,
	historicalMessages []openai.ChatCompletionMessageParamUnion) string {
	return singleChatCompletion(ctx, openaiClient, input, LookupConversationHistorySystemMessage, historicalMessages)
}

func singleChatCompletion(
	ctx context.Context,
	openaiClient openai.Client,
	input, systemMessage string,
	historicalMessages []openai.ChatCompletionMessageParamUnion) string {
	messages := []openai.ChatCompletionMessageParamUnion{
		{
			OfSystem: &openai.ChatCompletionSystemMessageParam{
				Content: openai.ChatCompletionSystemMessageParamContentUnion{
					OfString: openai.String(systemMessage),
				},
			},
		},
	}

	if len(historicalMessages) > 0 {
		messages = append(messages, historicalMessages...)
	}

	messages = append(messages, openai.ChatCompletionMessageParamUnion{
		OfUser: &openai.ChatCompletionUserMessageParam{
			Content: openai.ChatCompletionUserMessageParamContentUnion{
				OfString: openai.String(input),
			},
		},
	})

	chatResp, err := openaiClient.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Messages: messages,
			Model:    openai.ChatModelGPT4,
		})
	if err != nil {
		log.Fatalln("failed to chat response:", err)
	}

	return chatResp.Choices[0].Message.Content
}
