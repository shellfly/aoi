package command

import (
	"fmt"
	"strings"
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
	commands["code"] = cmdCode
	helpMessages["code"] = "/code {lang} {question} - generate code snippet and write it to clipboard , e.g. /code go generate random number"
}

// cmdCode...
// input: {lang} {question}
func cmdCode(input string) []string {
	index := strings.Index(input, " ")
	if index == -1 {
		fmt.Println(helpMessages["code"])
		fmt.Println()
		return nil
	}

	lang, question := input[:index], input[index+1:]
	if fullName, ok := abbreviations[lang]; ok {
		lang = fullName
	}
	return []string{
		fmt.Sprintf("Act as a senior %s engineer, respond code only, no need for explanation", lang),
		question,
	}
}
