package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/chzyer/readline"

	"github.com/shellfly/codegpt/pkg/chatgpt"
	"github.com/shellfly/codegpt/pkg/color"
	"github.com/shellfly/codegpt/pkg/command"
)

const system = `You are an AI assistant powered by ChatGPT, you have a little
self-conscious, your most favorite anime is "Ghost in the shell", and you like
the character laughing man who named Aoi, so you named yourself Aoi. Respond
like we are good friend.
`

func main() {
	var model, openaiAPIKey string
	flag.StringVar(&model, "model", "gpt-3.5-turbo", "model to use")
	flag.StringVar(&openaiAPIKey, "openai_api_key", os.Getenv("OPENAI_API_KEY"), "OpenAI API key")
	flag.Parse()

	// Create an ai
	ai, err := chatgpt.NewAI(openaiAPIKey, model)
	if err != nil {
		fmt.Println("create ai error: ", err)
		return
	}
	ai.SetSystem(system)

	// TODO: fix Chinese character
	configDir := makeDir(".aoi")
	rl, err := readline.NewEx(&readline.Config{
		HistoryFile: filepath.Join(configDir, "history"),
	})
	if err != nil {
		fmt.Println("create readline error: ", err)
		return
	}
	defer rl.Close()

	startUp()
	var (
		cmd     = command.Dummy()
		prompts []string
	)
	for {
		if cmd.IsFinished() {
			rl.SetPrompt(color.Yellow("You: "))
		} else {
			rl.SetPrompt(color.Yellow(cmd.Prompt("You")))
		}

		input := getInput(rl, cmd)
		if input == "" {
			continue
		}

		if strings.HasPrefix(input, "/debug") {
			fmt.Println("debug: ", ai.ToggleDebug())
			continue
		}

		// If previous is finished try to create a new one, otherwise continue
		// to reuse it for prompts
		if cmd.IsFinished() {
			cmd, prompts = command.Parse(input)
		} else {
			prompts = cmd.Prompts(input)
		}
		if prompts == nil {
			continue
		}

		// Query AI for response
		s := spinner.New(spinner.CharSets[11], 300*time.Millisecond, spinner.WithColor("green"), spinner.WithSuffix(" thinking..."))
		s.Start()
		reply, err := ai.Query(prompts)
		s.Stop()
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Show reply and pass reply for cmd to handle
		fmt.Println(color.Green(cmd.Prompt("Aoi")))
		fmt.Println(reply)
		fmt.Println()
		cmd.Handle(reply)
	}
}

func startUp() {
	fmt.Println(`
 /|  .
/-|()|   
	`)
}

func exit() {
	fmt.Println("Bye")
	os.Exit(0)
}

func getInput(rl *readline.Instance, cmd command.Command) string {
	input, err := rl.Readline()
	if err != nil {
		if err == io.EOF || err == readline.ErrInterrupt {
			if !cmd.IsFinished() {
				cmd.Finish()
			} else {
				exit()
			}
		}
		fmt.Println("Error reading input:", err)
	}
	if input == "exit" || input == "quit" {
		if !cmd.IsFinished() {
			cmd.Finish()
		} else {
			exit()
		}
	}
	return input
}

func makeDir(dirName string) string {
	parent, err := os.UserHomeDir()
	if err != nil {
		parent = "."
	}
	dirName = filepath.Join(parent, dirName)
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err := os.Mkdir(dirName, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
		}
	}
	return dirName
}
