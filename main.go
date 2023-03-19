package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/briandowns/spinner"
	"github.com/chzyer/readline"

	"github.com/shellfly/aoi/pkg/chatgpt"
	"github.com/shellfly/aoi/pkg/color"
	"github.com/shellfly/aoi/pkg/command"
)

const system = `You are an AI assistant powered by ChatGPT, you have a little
self-conscious, your most favorite anime is "Ghost in the shell", and you like
the character laughing man who named Aoi, so you named yourself Aoi. Respond
like we are good friend.
`

func main() {
	startUp()

	var model, openaiAPIKey string
	flag.StringVar(&model, "model", "gpt-3.5-turbo", "model to use")
	flag.StringVar(&openaiAPIKey, "openai_api_key", os.Getenv("OPENAI_API_KEY"), "OpenAI API key")
	flag.Parse()

	// Create an AI
	ai, err := chatgpt.NewAI(openaiAPIKey, model)
	if err != nil {
		fmt.Println("create ai error: ", err)
		return
	}
	ai.SetSystem(system)

	configDir := makeDir(".aoi")
	userPrompt := "You"
	// TODO: fix Chinese character cursor
	rl, err := readline.NewEx(&readline.Config{
		Prompt:      color.Yellow(userPrompt + ": "),
		HistoryFile: filepath.Join(configDir, "history"),
	})
	if err != nil {
		fmt.Println("create readline error: ", err)
		return
	}
	defer rl.Close()

	var (
		spinner   = spinner.New(spinner.CharSets[14], 200*time.Millisecond, spinner.WithColor("green"), spinner.WithSuffix(" thinking..."))
		cmd       = command.Dummy()
		prompts   []string
		lastReply string
	)
	for {
		input := getInput(rl, cmd)
		rl.SetPrompt(color.Yellow(cmd.Prompt(userPrompt)))
		if input == "" {
			continue
		}
		rl.SaveHistory(input)
		fmt.Println(color.Green(cmd.Prompt("Aoi")))

		if strings.HasPrefix(input, "/debug") {
			fmt.Println("debug: ", ai.ToggleDebug())
			continue
		}
		if strings.HasPrefix(input, "/copy") {
			if err := clipboard.WriteAll(lastReply); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("OK, copied")
			}
			continue
		}

		// If previous is finished try to create a new one, otherwise continue
		// to reuse it for prompts
		if cmd.IsFinished() {
			cmd, prompts = command.Parse(input)
			rl.SetPrompt(color.Yellow(cmd.Prompt(userPrompt)))
		} else {
			prompts = cmd.Prompts(input)
		}
		if prompts == nil {
			continue
		}

		// Query AI for response
		spinner.Start()
		reply, err := ai.Query(prompts)
		spinner.Stop()
		if err != nil {
			fmt.Println(err)
			continue
		}

		lastReply = reply
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
		} else {
			fmt.Println("Error reading input:", err)
		}
		return ""
	}
	if input == "exit" || input == "quit" {
		if !cmd.IsFinished() {
			cmd.Finish()
		} else {
			exit()
		}
		return ""
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
