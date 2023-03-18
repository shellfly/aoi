package command

import (
	"strings"
)

var commands = map[string]Command{}

type Command interface {
	/*
		required for a new command
	*/
	Name() string
	Help() string
	Run(string) []string // Run command and return optional ChatGPT prompts

	/*
		optional for new command, could inherit from dummyCommand
	*/
	Handle(string) // handle reply

	IsFinished() bool        // multiple commands mode
	Prompts(string) []string // multiple commands mode, continue generate prompts
	Close()                  // multiple commands mode, clean up

	Prompt(string) string // set custom terminal prompt
}

type dummyCommand struct{}

func (*dummyCommand) Name() string                  { return "dummy" }
func (*dummyCommand) Help() string                  { return "" }
func (*dummyCommand) Run(input string) []string     { return []string{input} }
func (*dummyCommand) Handle(string)                 {}
func (*dummyCommand) IsFinished() bool              { return true }
func (*dummyCommand) Prompts(input string) []string { return []string{input} }
func (*dummyCommand) Close()                        {}
func (*dummyCommand) Prompt(p string) string        { return p + ": " }

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
	return cmd, cmd.Run(input)
}
