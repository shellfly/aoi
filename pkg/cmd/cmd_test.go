package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	system, prompt := Parse("/py generate number in 1 to 100")
	assert.Equal(t, engSystem("py"), system)
	assert.Equal(t, "generate number in 1 to 100", prompt)
}
