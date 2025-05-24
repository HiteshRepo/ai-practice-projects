package utils

import (
	"fmt"
	"react/models"
	"regexp"
	"strings"
)

var (
	ErrNoActionFound = fmt.Errorf("action not found in assistant response")
)

func IsEmpty(arg string) bool {
	return len(arg) == 0 ||
		len(strings.TrimSpace(arg)) == 0 ||
		strings.Contains(strings.ToUpper(strings.TrimSpace(arg)), "NONE")
}

const actionRegex = `Action:\s*([^:]+):\s*(.+)`

func ActionExtractor(
	input string,
) (models.Action, error) {
	actionRegex := regexp.MustCompile(actionRegex)
	matches := actionRegex.FindStringSubmatch(input)

	if len(matches) == 0 {
		return models.Action{}, ErrNoActionFound
	}

	if len(matches) < 3 {
		return models.Action{}, fmt.Errorf("incorrect format of action in assistant response")
	}

	functionName := strings.TrimSpace(matches[1])
	arguments := make([]string, 0)

	for i := 2; i < len(matches); i++ {
		arguments = append(arguments, strings.TrimSpace(matches[i]))
	}

	return models.Action{
		FunctionName: functionName,
		Arguments:    arguments,
	}, nil
}
