package command

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/atotto/clipboard"
)

var abbreviations = map[string]string{
	"go": "golang",
	"js": "javascript",
	"oc": "objective-c",
	"pg": "postgres",
	"py": "python",
	"rn": "react native",
	"re": "regular expression",
}

func init() {
	commands["code"] = &Code{}
}

type Code struct {
	dummyCommand
}

func (c *Code) Name() string {
	return "code"
}

func (c *Code) Help() string {
	return "/code {lang} {question} - Generate code snippets and copy to the clipboard , e.g. /code go generate random number"
}

// Prompts expand input like "{lang} {question}" to code generation prompts
func (c *Code) Prompts(input string) []string {
	index := strings.Index(input, " ")
	if index == -1 {
		fmt.Println(c.Help())
		return nil
	}

	lang, question := input[:index], input[index+1:]
	if fullName, ok := abbreviations[lang]; ok {
		lang = fullName
	}
	return []string{
		fmt.Sprintf("You are a programming expert on %s, respond code only, don't add comment, don't do explanation", lang),
		question,
	}
}

// Handle copy code in the reply to clipboard, and return the original reply
func (c *Code) Handle(reply string) {
	code := reply
	if strings.Contains(reply, "```") {
		code = extractCode(reply)
	}
	fmt.Println(code)
	fmt.Println()
	if code != "" {
		if err := clipboard.WriteAll(code); err != nil {
			fmt.Printf("failed to copy to clipboard: %v", err)
		}
	}
}

// extractCode extract first markdown code snippet in text
func extractCode(text string) string {
	//sre := regexp.MustCompile("(?sm)```.*\n(.*)\n```.*")
	re := regexp.MustCompile("(?sm)^```" + `\w*(.*?)` + "```$")
	matches := re.FindStringSubmatch(text)
	if len(matches) > 0 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}
