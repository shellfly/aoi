package command

import (
	"fmt"
	"sort"
)

func init() {
	commands["help"] = cmdHelp
	helpMessages["help"] = "/help -- show the help message"
}

func cmdHelp(input string) []string {
	cmds := make([]string, 0, len(helpMessages))
	for cmd := range helpMessages {
		cmds = append(cmds, cmd)
	}

	sort.Strings(cmds)
	for _, cmd := range cmds {
		fmt.Println(helpMessages[cmd])
	}
	return nil
}
