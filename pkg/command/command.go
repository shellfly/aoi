package command

import (
	"strings"
)

type commandFunc func(string) []string

var (
	commands     = map[string]commandFunc{}
	helpMessages = map[string]string{}
)

// Parse parse slash command in input and generate prompts for ChatGPT
func Parse(input string) []string {
	if !strings.HasPrefix(input, "/") {
		return []string{input}
	}

	input = input[1:]
	index := strings.Index(input, " ")
	var cmd string
	if index == -1 {
		cmd, input = input, ""
	} else {
		cmd, input = input[:index], input[index+1:]
	}

	cmdFunc, ok := commands[cmd]
	if !ok {
		cmdFunc = cmdHelp
	}
	return cmdFunc(input)
}
