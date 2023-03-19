package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	cmd, prompts := Parse("/code py generate number in 1 to 100")
	assert.Equal(t, "code", cmd.Name())
	assert.Equal(t, "generate number in 1 to 100", prompts[1])

	cmd, _ = Parse("/shell generate number in 1 to 100")
	assert.Equal(t, "shell", cmd.Name())
}
