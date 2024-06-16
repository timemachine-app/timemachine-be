package util

import "strings"

func CleanOpenAIJson(input string) string {
	// Trim the prefix "json" or "```" from the start of the string
	if strings.HasPrefix(input, "json") {
		input = strings.TrimPrefix(input, "json")
	} else if strings.HasPrefix(input, "```") {
		input = strings.TrimPrefix(input, "```")
	}

	// Trim the postfix "```" from the end of the string
	if strings.HasSuffix(input, "```") {
		input = strings.TrimSuffix(input, "```")
	}

	// Trim any leading or trailing whitespace
	input = strings.TrimSpace(input)

	return input
}
