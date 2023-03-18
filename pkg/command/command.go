package command

import (
	"strings"
)

var commands = map[string]Command{}

type Command interface {
	Name() string
	Help() string

	// expand input to prompts
	Expand(string) []string
	// handle reply
	Handle(string)
}

type dummyHandler struct{}

func (dummyHandler) Handle(string) {}

// Parse parse slash command in input and generate prompts for ChatGPT
func Parse(input string) (cmd Command, prompts []string) {
	if !strings.HasPrefix(input, "/") {
		return nil, []string{input}
	}

	input = input[1:]
	index := strings.Index(input, " ")
	var cmdName string
	if index == -1 {
		cmdName, input = input, ""
	} else {
		cmdName, input = input[:index], input[index+1:]
	}

	cmd, ok := commands[cmdName]
	if !ok {
		cmd = commands["help"]
	}
	return cmd, cmd.Expand(input)
}
