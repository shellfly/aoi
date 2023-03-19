package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractCod(t *testing.T) {
	sql := "SELECT name FROM sqlite_master WHERE type='table' ORDER BY name;"
	text := `
	To list tables, you can execute the following SQL statement:
	` + wrapCode("", sql) + `
	This will return a list of all table names in the database.
	`
	code := extractCode(text)
	assert.Equal(t, sql, code)

	shell := wrapCode("bash", "ls")
	code = extractCode(shell)
	assert.Equal(t, "ls", code)
}

func wrapCode(lang string, code string) string {
	return "```" + lang + `
` + code + `
` + "```"
}
