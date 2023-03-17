package cmd

import (
	"fmt"
	"strings"
)

// Parse parse slash command in input and generate "system" and "prompt" string
// for ChatGPT
func Parse(input string) (system, prompt string) {
	if !strings.HasPrefix(input, "/") {
		return "", input
	}
	input = input[1:]
	index := strings.Index(input, " ")
	cmd := input[:index]
	prompt = input[index+1:]

	return engSystem(cmd), prompt
}

// engSystem expand the language to a engineer system instruction
func engSystem(lang string) string {
	const format = "Act as a senior %s engineer, respond code only, no need for explanation"
	if language, ok := abbreviations[lang]; ok {
		lang = language
	}
	return fmt.Sprintf(format, lang)
}
