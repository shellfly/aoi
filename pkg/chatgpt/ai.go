package chatgpt

import (
	"context"
	"fmt"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

const MessageLimit = 100

type AI struct {
	client *openai.Client

	system   string
	model    string
	messages []openai.ChatCompletionMessage

	debug bool
}

func NewAI(client *openai.Client, model string) *AI {
	messages := make([]openai.ChatCompletionMessage, 0, 2*MessageLimit)
	ai := &AI{
		client:   client,
		model:    model,
		messages: messages,
		debug:    false,
	}
	return ai
}

func (ai *AI) SetSystem(system string) {
	ai.system = system
	ai.messages = []openai.ChatCompletionMessage{NewMessage(openai.ChatMessageRoleSystem, system)}
}

// limitTokens make sure messages are under tokens limit
// TODO: accurate way to control tokens limit
// https://help.openai.com/en/articles/4936856-what-are-tokens-and-how-to-count-them
func (ai *AI) limitTokens() {
	if len(ai.messages) < MessageLimit {
		return
	}

	// keep last MessageLimit messages and the system message
	copy(ai.messages[1:], ai.messages[len(ai.messages)-MessageLimit:])
	ai.messages = ai.messages[:MessageLimit]
}

func (ai *AI) Query(prompts []string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	for _, prompt := range prompts {
		ai.messages = append(ai.messages, NewMessage(openai.ChatMessageRoleUser, prompt))
	}
	ai.limitTokens()

	if ai.debug {
		fmt.Println("---debug---")
		for _, msg := range ai.messages {
			fmt.Println(msg)
		}
		fmt.Println("---debug---")
	}
	// Set the request parameters for the completion API
	req := openai.ChatCompletionRequest{
		Model:    ai.model,
		Messages: ai.messages,
	}

	// Send the completion API request to OpenAI and get the response
	resp, err := ai.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	reply := resp.Choices[0].Message.Content
	ai.messages = append(ai.messages, NewMessage(openai.ChatMessageRoleAssistant, reply))
	return strings.TrimSpace(reply), nil
}

func (ai *AI) ToggleDebug() bool {
	ai.debug = !ai.debug
	if ai.debug {
		fmt.Println(ai.messages)
	}
	return ai.debug
}

func (ai *AI) Reset() {
	ai.messages = ai.messages[:1]
}

func NewMessage(role, text string) openai.ChatCompletionMessage {
	return openai.ChatCompletionMessage{
		Role:    role,
		Content: text,
	}
}
