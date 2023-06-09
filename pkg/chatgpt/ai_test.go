package chatgpt

import (
	"fmt"
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
)

func TestAI(t *testing.T) {
	client := openai.NewClient("api key")
	ai := NewAI(client, "model")
	t.Run("limit tokens", func(t *testing.T) {
		ai.messages = make([]openai.ChatCompletionMessage, MessageLimit+2)
		ai.messages[0] = NewMessage("system", "message")
		for i := 1; i < MessageLimit+1; i++ {
			ai.messages[i] = NewMessage("user", fmt.Sprintf("message %d", i))
		}
		ai.limitTokens()
		assert.Equal(t, MessageLimit, len(ai.messages))
		assert.Equal(t, "system", ai.messages[0].Role)
		assert.Equal(t, "message 2", ai.messages[1].Content)
	})
}
