package utils

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	ErrNoActionFound = fmt.Errorf("action not found in assistant response")
)

type Action struct {
	FunctionName string
	Arguments    []string
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
