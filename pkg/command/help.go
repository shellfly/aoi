package command

import (
	"fmt"
	"sort"
)

func init() {
	commands["help"] = &Help{}
}

type Help struct {
	dummyCommand
}

func (c *Help) Name() string {
	return "help"
}
func (c *Help) Help() string {
	return "/help - show the help message"
}

func (c *Help) Run(input string) []string {
	names := make([]string, 0, len(commands))
	for cmd := range commands {
		names = append(names, cmd)
	}

	sort.Strings(names)
	for _, name := range names {
		fmt.Println(commands[name].Help())
	}
	fmt.Println()
	return nil
}
