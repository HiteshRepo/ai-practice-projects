package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/openai/openai-go"
)

var (
	ErrNoActionFound = fmt.Errorf("action not found in assistant response")
)

type Action struct {
	FunctionName string
	Arguments    []string
	ToolCallID   string
}

func IsEmpty(arg string) bool {
	return arg == "{}" ||
		len(arg) == 0 ||
		len(strings.TrimSpace(arg)) == 0 ||
		strings.ToUpper(strings.TrimSpace(arg)) == `"NONE"`
}

const actionRegex = `Action:\s*([^:]+):\s*(.+)`

func ActionExtractor(
	input string,
) (Action, error) {
	actionRegex := regexp.MustCompile(actionRegex)
	matches := actionRegex.FindStringSubmatch(input)

	if len(matches) == 0 {
		return Action{}, ErrNoActionFound
	}

	if len(matches) < 3 {
		return Action{}, fmt.Errorf("incorrect format of action in assistant response")
	}

	functionName := strings.TrimSpace(matches[1])
	arguments := make([]string, 0)

	for i := 2; i < len(matches); i++ {
		arguments = append(arguments, strings.TrimSpace(matches[i]))
	}

	return Action{
		FunctionName: functionName,
		Arguments:    arguments,
	}, nil
}

func ActionsFromResponseToolCalls(toolCalls []openai.ChatCompletionMessageToolCall) []Action {
	actions := make([]Action, 0)

	for _, tool := range toolCalls {
		action := Action{
			FunctionName: tool.Function.Name,
			Arguments:    []string{tool.Function.Arguments},
			ToolCallID:   tool.ID,
		}

		actions = append(actions, action)
	}

	return actions
}
