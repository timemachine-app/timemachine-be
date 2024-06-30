package util

import "strings"

func CleanLLMJson(input string) string {
	// Trim the prefix "json" and "```" from the start of the string
	if strings.HasPrefix(input, "```") {
		input = strings.TrimPrefix(input, "```")
	}
	if strings.HasPrefix(input, "json") {
		input = strings.TrimPrefix(input, "json")
	}

	// Trim the postfix "```" from the end of the string
	if strings.HasSuffix(input, "```") {
		input = strings.TrimSuffix(input, "```")
	}

	// Trim any leading or trailing whitespace
	input = strings.TrimSpace(input)

	return input
}
