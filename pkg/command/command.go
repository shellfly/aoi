package command

import (
	"fmt"
	"strings"
)

var commands = map[string]Command{}

type Command interface {
	Name() string
	Help() string
	Init(string) string      // initialize command by input
	Prompts(string) []string // return optional ChatGPT prompts

	Handle(string) // handle reply

	IsFinished() bool // multiple commands mode
	Finish()          // multiple commands mode, clean up

	Prompt(string) string // set custom terminal prompt
}

type dummyCommand struct{}

func (*dummyCommand) Name() string                  { return "" }
func (*dummyCommand) Help() string                  { return "" }
func (*dummyCommand) Init(input string) string      { return input }
func (*dummyCommand) Prompts(input string) []string { return []string{input} }
func (*dummyCommand) Handle(reply string) {
	fmt.Println(reply)
	fmt.Println()
}
func (*dummyCommand) IsFinished() bool       { return true }
func (*dummyCommand) Finish()                {}
func (*dummyCommand) Prompt(p string) string { return p + ": " }

// Dummy return the dummy command
func Dummy() Command {
	return &dummyCommand{}
}

// Parse parse slash command in input and generate prompts for ChatGPT
func Parse(input string) (cmd Command, prompts []string) {
	if !strings.HasPrefix(input, "/") {
		return Dummy(), []string{input}
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
	input = cmd.Init(input)
	if input != "" {
		prompts = cmd.Prompts(input)
	}
	return cmd, prompts
}
