package openai

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
)

func CreateAssistant(
	ctx context.Context,
	openaiClient openai.Client,
	instructions string,
	name string,
	tools []openai.AssistantToolUnionParam,
	toolResources *openai.BetaAssistantNewParamsToolResources,
) (*openai.Assistant, error) {
	params := openai.BetaAssistantNewParams{
		Model:        "gpt-4",
		Name:         openai.String(name),
		Instructions: openai.String(instructions),
	}

	if len(tools) > 0 {
		params.Tools = tools
	}

	if toolResources != nil {
		params.ToolResources = *toolResources
	}

	return openaiClient.Beta.Assistants.New(ctx, params)
}

func IsAssistantCreated(
	ctx context.Context,
	client openai.Client,
	name string) (bool, *openai.Assistant, error) {
	createdAssistants, err := client.Beta.Assistants.List(
		ctx,
		openai.BetaAssistantListParams{})
	if err != nil {
		return false, nil, err
	}

	for _, assistant := range createdAssistants.Data {
		if assistant.Name == name {
			return true, &assistant, nil
		}
	}

	after := ""
	if len(createdAssistants.Data) > 0 {
		after = createdAssistants.Data[len(createdAssistants.Data)-1].ID
	}

	for {
		if !createdAssistants.HasMore {
			break
		}

		nextPage, err := client.Beta.Assistants.List(
			ctx,
			openai.BetaAssistantListParams{
				After: param.NewOpt(after),
			})
		if err != nil {
			return false, nil, err
		}

		if len(nextPage.Data) == 0 {
			break
		}

		for _, assistant := range nextPage.Data {
			if assistant.Name == name {
				return true, &assistant, nil
			}
		}

		after = nextPage.Data[len(nextPage.Data)-1].ID
		createdAssistants = nextPage
	}

	return false, nil, nil
}
