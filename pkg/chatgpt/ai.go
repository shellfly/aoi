package chatgpt

import (
	"context"
	"errors"
	"fmt"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

type AI struct {
	client *openai.Client

	model    string
	messages []openai.ChatCompletionMessage

	debug bool
}

func NewAI(apiKey, model string) (*AI, error) {
	if apiKey == "" {
		return nil, errors.New("Please set the OPENAI_API_KEY environment variable")
	}

	// Create a new OpenAI API client with the provided API key
	client := openai.NewClient(apiKey)
	messages := []openai.ChatCompletionMessage{}
	ai := &AI{
		client:   client,
		model:    model,
		messages: messages,
		debug:    false,
	}
	return ai, nil
}

func (ai *AI) Query(system string, prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if system != "" {
		ai.messages = []openai.ChatCompletionMessage{NewMessage(openai.ChatMessageRoleSystem, system)}
	}

	ai.messages = append(ai.messages, NewMessage(openai.ChatMessageRoleUser, prompt))
	// TODO: accurate way to control tokens limit
	// https://help.openai.com/en/articles/4936856-what-are-tokens-and-how-to-count-them
	if len(ai.messages) > 100 {
		ai.messages = ai.messages[len(ai.messages)-100:]
	}

	if ai.debug {
		fmt.Println(ai.messages)
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
	return reply, nil
}

func (ai *AI) ClearMessages() {
	ai.messages = []openai.ChatCompletionMessage{}
}

func (ai *AI) ToggleDebugMode() {
	ai.debug = !ai.debug
}

func NewMessage(role, text string) openai.ChatCompletionMessage {
	return openai.ChatCompletionMessage{
		Role:    role,
		Content: text,
	}
}
