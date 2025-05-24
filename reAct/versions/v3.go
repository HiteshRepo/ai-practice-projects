package versions

import (
	"context"

	"log"
	"react/constants"
	"react/tools"
	"react/utils"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/shared"
	"github.com/openai/openai-go/shared/constant"
)

var toolsParams = []openai.ChatCompletionToolParam{
	{
		Type: constant.Function(openai.AssistantToolChoiceTypeFunction),
		Function: shared.FunctionDefinitionParam{
			Name:        "getCurrentWeather",
			Description: openai.String("Get the current weather"),
			Parameters: map[string]any{
				"type":       "object",
				"properties": tools.GetCurrentWeatherAllProperties,
				"required":   tools.GetCurrentWeatherRequiredProperties,
			},
		},
	},
	{
		Type: constant.Function(openai.AssistantToolChoiceTypeFunction),
		Function: shared.FunctionDefinitionParam{
			Name:        "getLocation",
			Description: openai.String("Get the user's current location"),
			Parameters: map[string]any{
				"type":       "object",
				"properties": tools.GetLocationAllProperties,
				"required":   tools.GetLocationRequiredProperties,
			},
		},
	},
}

func V3(
	ctx context.Context,
	openaiClient openai.Client,
	wsClient *utils.WeatherStack,
	query string) {
	messages := []openai.ChatCompletionMessageParamUnion{}

	systemMessage := openai.ChatCompletionMessageParamUnion{
		OfSystem: &openai.ChatCompletionSystemMessageParam{
			Content: openai.ChatCompletionSystemMessageParamContentUnion{
				OfString: openai.String(constants.BriefReActSystemPrompt),
			},
		},
	}

	userMessage := openai.ChatCompletionMessageParamUnion{
		OfUser: &openai.ChatCompletionUserMessageParam{
			Content: openai.ChatCompletionUserMessageParamContentUnion{
				OfString: openai.String(query),
			},
		},
	}

	messages = append(messages, systemMessage)
	messages = append(messages, userMessage)

	for i := 0; i < constants.MaxIterations; i++ {
		log.Printf("Iteration #%d", i+1)

		resp, err := openaiClient.Chat.Completions.New(
			ctx,
			openai.ChatCompletionNewParams{
				Model:    "gpt-4",
				Messages: messages,
				// https://platform.openai.com/docs/api-reference/chat/create#chat-create-tools
				Tools: toolsParams,
			},
		)
		if err != nil {
			log.Fatalln(err)
		}

		messages = append(messages, resp.Choices[0].Message.ToParam())

		switch resp.Choices[0].FinishReason {
		case "tool_calls":
			actions, err := tools.ActionsFromResponseToolCalls(resp.Choices[0].Message.ToolCalls)
			if err != nil {
				log.Fatalln(err)
			}

			for _, action := range actions {
				actionResponse := tools.InvokeValidAction(wsClient, action)

				toolMessage := openai.ChatCompletionMessageParamUnion{
					OfTool: &openai.ChatCompletionToolMessageParam{
						ToolCallID: action.ToolCallID,
						Content: openai.ChatCompletionToolMessageParamContentUnion{
							OfString: openai.String(actionResponse),
						},
					},
				}

				messages = append(messages, toolMessage)
			}

		case "stop":
			log.Println(resp.Choices[0].Message.Content)
			return
		}
	}
}
