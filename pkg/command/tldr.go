package command

import (
	"fmt"
)

func init() {
	commands["tldr"] = &TLDR{}
}

type TLDR struct {
	dummyCommand
}

func (c *TLDR) Name() string {
	return "TLDR"
}

func (c *TLDR) Help() string {
	return "/tldr {command} - Get a tl;dr explanation of a command"
}

// Prompts generate tldr prompt
func (c *TLDR) Prompts(input string) []string {
	return []string{
		fmt.Sprintf(`
I want you to act as man pages. I will ask you a command and you will show simpler doc and practical examples, don't do explanation.
The command is: %s`, input),
	}
}
