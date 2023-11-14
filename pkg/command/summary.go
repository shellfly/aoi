package command

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func init() {
	commands["summary"] = &Summary{}
}

type Summary struct {
	dummyCommand
}

func (c *Summary) Name() string {
	return "summary"
}

func (c *Summary) Help() string {
	return "/summary - generate summary from a URL"
}

// Prompts expand input like "{lang} {question}" to Summary generation prompts
func (c *Summary) Prompts(input string) []string {
	parts := strings.Split(input, " ")
	var lang, url string
	if len(parts) == 2 {
		lang = parts[0]
		url = parts[1]
	} else {
		url = input
	}

	content := crawl(url)
	prompt := fmt.Sprintf(`Generate a summary of the below text content.n\nText:"""\n%s\n"""`, content)
	if language, ok := Languages[lang]; ok {
		prompt = prompt + "\nTranslate the response to " + language
	}
	return []string{prompt}
}

func crawl(url_input string) string {
	// Validate the URL
	if _, err := url.ParseRequestURI(url_input); err != nil {
		fmt.Println("Invalid URL: ", err)
		return ""
	}

	resp, err := http.Get(url_input)
	if err != nil {
		fmt.Println("Error fetching URL: ", err)
		return ""
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body: ", err)
		return ""
	}
	doc.Find("nav, script, iframe, style, footer").Remove()
	r := regexp.MustCompile(`\s+`)
	text := doc.Find("body").Text()
	text = r.ReplaceAllString(text, " ")
	return text
}
