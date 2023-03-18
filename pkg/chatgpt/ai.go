package chatgpt

import (
	"context"
	"errors"
	"fmt"
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

func NewAI(apiKey, model string) (*AI, error) {
	if apiKey == "" {
		return nil, errors.New("Please set the OPENAI_API_KEY environment variable")
	}

	// Create a new OpenAI API client with the provided API key
	client := openai.NewClient(apiKey)
	messages := make([]openai.ChatCompletionMessage, 2*MessageLimit)
	ai := &AI{
		client:   client,
		model:    model,
		messages: messages,
		debug:    false,
	}
	return ai, nil
}
func (ai *AI) SetSystem(system string) {
	ai.system = system
	ai.messages = []openai.ChatCompletionMessage{NewMessage(openai.ChatMessageRoleSystem, system)}
}

func (ai *AI) Query(prompts []string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	for _, prompt := range prompts {
		ai.messages = append(ai.messages, NewMessage(openai.ChatMessageRoleUser, prompt))
	}
	// TODO: accurate way to control tokens limit
	// https://help.openai.com/en/articles/4936856-what-are-tokens-and-how-to-count-them
	if len(ai.messages) > MessageLimit {
		// limit to 100 messages and always keep the system message
		copy(ai.messages[1:], ai.messages[len(ai.messages)-MessageLimit:])
		ai.messages = ai.messages[:MessageLimit]
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

func (ai *AI) ToggleDebug() bool {
	ai.debug = !ai.debug
	return ai.debug
}

func NewMessage(role, text string) openai.ChatCompletionMessage {
	return openai.ChatCompletionMessage{
		Role:    role,
		Content: text,
	}
}
