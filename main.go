package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/chzyer/readline"

	"github.com/shellfly/codegpt/pkg/chatgpt"
	"github.com/shellfly/codegpt/pkg/cmd"
	"github.com/shellfly/codegpt/pkg/color"
)

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

	// Create a new readline instance to read user input from the console
	// TODO: fix Chinese character
	rl, err := readline.NewEx(&readline.Config{
		Prompt:      color.Yellow("You: "),
		HistoryFile: "/tmp/.codegpt-readline.tmp",
	})
	if err != nil {
		fmt.Println("create readline error: ", err)
		return
	}
	defer rl.Close()

	startUp()
	for {
		// Ask the user for input
		input, err := rl.Readline()
		if err != nil {
			if err == io.EOF || err == readline.ErrInterrupt {
				exit()
			}
			fmt.Println("Error reading input:", err)
			continue
		}
		// If the user entered an empty input, skip to the next iteration of the loop
		if input == "" {
			continue
		}

		// If the user entered the "exit" command, break out of the loop and exit the program
		if input == "exit" || input == "quit" {
			exit()
		}

		if strings.HasPrefix(input, "/debug") {
			ai.ToggleDebugMode()
			continue
		}

		system, prompt := cmd.Parse(input)
		reply, err := ai.Query(system, prompt)
		if err != nil {
			fmt.Println("Oops, ", err)
			continue
		}
		reply = strings.TrimSpace(reply)
		err = copyCode(reply)
		if err != nil {
			fmt.Printf("failed to copy to clipboard: %v", err)
		}

		fmt.Println(color.Green("AI:"))
		fmt.Println(reply)
	}
}

// extractCode extract first markdown code snippet in text
func extractCode(text string) string {
	re := regexp.MustCompile("(?sm)^```" + ` ?\w*(.*?)` + "```$")
	matches := re.FindStringSubmatch(text)
	if len(matches) > 0 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

// copyCode copy the first code snippet in text to clipboard
func copyCode(text string) error {
	code := extractCode(text)
	if code != "" {
		return clipboard.WriteAll(code)
	}
	return nil
}

func startUp() {
	fmt.Println(`
_________            .___       _____________________________
\_   ___ \  ____   __| _/____  /  _____/\______   \__    ___/
/    \  \/ /  _ \ / __ |/ __ \/   \  ___ |     ___/ |    |   
\     \___(  <_> ) /_/ \  ___/\    \_\  \|    |     |    |   
\________ /\____/\____ |\____>\_________/|____|     |____|   
	`)
}

func exit() {
	fmt.Println("Bye")
	os.Exit(0)
}
