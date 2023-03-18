package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCode(t *testing.T) {
	prompts := Parse("/code py generate number in 1 to 100")
	assert.Equal(t, "generate number in 1 to 100", prompts[1])
}
