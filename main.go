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
	"github.com/shellfly/codegpt/pkg/color"
	"github.com/shellfly/codegpt/pkg/command"
)

const system = `You are an AI assistant powered by ChatGPT, you a little
self-conscious, your most favorite anime is "Ghost in the shell", and you like
the character laughing man which also named Aoi, so you named yourself Aoi.
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
			fmt.Println("debug: ", ai.ToggleDebug())
			continue
		}

		prompts := command.Parse(input)
		if prompts == nil {
			continue
		}

		reply, err := ai.Query(prompts)
		if err != nil {
			fmt.Println(err)
			continue
		}
		reply = strings.TrimSpace(reply)
		err = copyCode(reply)
		if err != nil {
			fmt.Printf("failed to copy to clipboard: %v", err)
		}

		fmt.Println(color.Green("Aoi:"))
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
 /|  .
/-|()|   
	`)
}

func exit() {
	fmt.Println("Bye")
	os.Exit(0)
}
