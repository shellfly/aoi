package command

import (
	"fmt"
	"strings"
)

// https://en.wikipedia.org/wiki/IETF_language_tag
var Languages = map[string]string{
	"ar":      "Arabic",
	"cn":      "Chinese",
	"de":      "German",
	"en":      "English",
	"fr":      "French",
	"hi":      "Hindi",
	"ja":      "Japanese",
	"jp":      "Japanese",
	"pt":      "Portuguese",
	"ru":      "Russian",
	"spa":     "Spanish",
	"es":      "Spanish",
	"zh-hant": "Traditional Chinese",
	"zh-tw":   "Traditional Chinese",
	"zh":      "Chinese",
}

func init() {
	commands["trans"] = &Trans{}
}

type Trans struct {
	dummyCommand
}

func (c *Trans) Name() string {
	return "translate"
}

func (c *Trans) Help() string {
	return "/trans {lang code} {text} - Translate text to a specified language"
}

// Run expand input like "{lang} {question}" to code generation prompts
func (c *Trans) Prompts(input string) []string {
	index := strings.Index(input, " ")
	if index == -1 {
		fmt.Println(c.Help())
		return nil
	}

	lang, text := input[:index], input[index+1:]
	if fullName, ok := Languages[lang]; ok {
		lang = fullName
	}
	return []string{
		fmt.Sprintf("Translate to %s: %s", lang, text),
	}
}
