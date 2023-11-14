package command

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
)

const (
	WINDOWS_OS = "windows"
	LINUX_OS   = "linux"
	MAC_OS     = "darwin"
)

func init() {
	commands["shell"] = &Shell{}
}

type Shell struct {
	dummyCommand
}

func (c *Shell) Name() string {
	return "shell"
}

func (c *Shell) Help() string {
	return "/shell - generate shell command and execute it"
}

// Prompts expand input like "{lang} {question}" to Shell generation prompts
func (c *Shell) Prompts(input string) []string {
	if strings.HasPrefix(input, ":") {
		c.Handle(input[1:])
		return nil
	}

	return []string{
		fmt.Sprintf(`
I want you to act as a terminal. I will ask you a question and you will reply with one-line command to do it.
I want you to only reply with the code, and nothing else. do not write explanations.
My question is how to %s on %s?
		`, input, getOSInfo()),
	}
}

// Handle execute shell command
func (c *Shell) Handle(reply string) {
	ExecShell(reply)
}

func ExecShell(command string) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	defer signal.Reset(os.Interrupt)

	command = strings.ReplaceAll(command, "`", "")
	fmt.Println(command)
	fmt.Println()
	var cmd *exec.Cmd
	if runtime.GOOS == WINDOWS_OS {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("bash", "-c", command)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return
	}
	done := make(chan struct{})
	go func() {
		_ = cmd.Wait()
		close(done)
	}()

	// Wait for the command to finish or for a signal to be received
	select {
	case <-done:
		return
	case <-sigChan:
		if runtime.GOOS == WINDOWS_OS {
			_ = cmd.Process.Signal(os.Kill)
		} else {
			_ = syscall.Kill(cmd.Process.Pid, syscall.SIGINT)
		}
		fmt.Println()
	}
}

func getOSInfo() string {
	var version string
	switch runtime.GOOS {
	case MAC_OS:
		version = darwinVersion()
	case LINUX_OS:
		version = linuxVersion()
	case WINDOWS_OS:
		version = windowsVersion()
	}
	return fmt.Sprintf("%s %s", runtime.GOOS, version)
}

func darwinVersion() string {
	cmd := exec.Command("sw_vers", "-productVersion")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	return string(output)
}

func linuxVersion() string {
	cmd := exec.Command("lsb_release", "-d", "-s")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(output)
}

func windowsVersion() string {
	cmd := exec.Command("systeminfo", "/fo", "csv", "/nh")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	reader := csv.NewReader(strings.NewReader(string(output)))
	records, err := reader.ReadAll()
	if err != nil {
		return ""
	}

	for _, record := range records {
		if len(record) >= 2 && record[0] == "OS Version:" {
			return record[1]
		}
	}
	return ""
}
