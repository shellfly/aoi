package command

import "testing"

func TestShell(t *testing.T) {
	cmd := "ls | sleep 1"
	ExecShell(cmd)
}
