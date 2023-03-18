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

// Run expand input like "{lang} {question}" to Shell generation prompts
func (c *Shell) Run(input string) []string {
	if strings.HasPrefix(input, ":") {
		fmt.Println(input[1:])
		c.Handle(input[1:])
		return nil
	}

	return []string{
		fmt.Sprintf(`
I want you to act as a terminal. I will ask you a question and you will reply with one-line command to do it, avoid pipeline if possible.
I want you to only reply with the code, and nothing else. do not write explanations.
My question is how to %s on %s?
		`, input, getOSInfo()),
	}
}

// Handle execute shell command
func (c *Shell) Handle(reply string) {
	ExecCommand(reply)
}

func ExecCommand(command string) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	defer signal.Reset(os.Interrupt)

	command = strings.ReplaceAll(command, "`", "")
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("bash", "-c", command)
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}

	output := make(chan []byte, 1)
	go func() {
		// TODO: real-time output
		out, _ := cmd.CombinedOutput()
		output <- out
	}()

	// Wait for the command to finish or for a signal to be received
	select {
	case out := <-output:
		fmt.Println(string(out))
		fmt.Println()
	case <-sigChan:
		var err error
		if runtime.GOOS == "windows" {
			err = cmd.Process.Signal(os.Kill)
		} else {
			err = syscall.Kill(-cmd.Process.Pid, syscall.SIGINT)
		}
		if err != nil {
			fmt.Println("stop process error: ", err)
		}
	}
}

func getOSInfo() string {
	var version string
	switch runtime.GOOS {
	case "darwin":
		version = darwinVersion()
	case "linux":
		version = linuxVersion()
	case "windows":
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
