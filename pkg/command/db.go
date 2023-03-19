package command

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
)

func init() {
	commands["db"] = &DB{}
}

// TODO: Leverage other tools to keep things simple
// use Go code to access various database
var clients = map[string]string{
	"mysql":    "mysql",
	"postgres": "psql",
	"sqlite":   "sqlite3",
}

type DB struct {
	dummyCommand

	dbType      string
	url         string
	client      string
	isFinished  bool
	initialized bool
}

func (c *DB) Name() string {
	return "database"
}

func (c *DB) Help() string {
	return "/db - generate SQL and execute it on database, e.g. /db {url} show tables"
}

// Prompt set terminal prompt for ssh command
func (c *DB) Prompt(p string) string {
	if c.isFinished {
		return p + ": "
	}
	return fmt.Sprintf("(/db %s) %s: ", c.dbType, p)
}

// Init...
func (c *DB) Init(input string) string {
	c.isFinished = true
	index := strings.Index(input, " ")
	if index == -1 {
		if err := c.setUrl(input); err != nil {
			fmt.Println(err)
			return ""
		}
		c.isFinished = false
		return ""
	}

	url, input := input[:index], input[index+1:]
	if err := c.setUrl(url); err != nil {
		fmt.Println(err)
		return ""
	}
	return input
}

func (c *DB) IsFinished() bool {
	return c.isFinished
}

func (c *DB) Prompts(input string) []string {
	if strings.HasPrefix(input, ":") {
		fmt.Println(input[1:])
		c.Handle(input[1:])
		return nil
	}

	if c.initialized {
		return []string{input}
	}
	c.initialized = true

	prompts := []string{
		fmt.Sprintf("Given these table definitions: \n %s\n", c.FetchTables()),
		fmt.Sprintf("You are a %s expert, reply with the code, and nothing else. do not write explanations.", c.dbType),
		input,
	}
	return prompts
}

func (c *DB) Finish() {
	c.dbType = ""
	c.url = ""
	c.isFinished = true
}

// Handle execute SQL
func (c *DB) Handle(reply string) {
	sql := reply
	if strings.Contains(reply, "```") {
		sql = extractCode(reply)
	}
	output, err := c.ExecSQL(sql)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(output)
	}
	fmt.Println()
}

func (c *DB) setUrl(url string) error {
	parts := strings.Split(url, "://")
	dbType, url := parts[0], parts[1]
	client, ok := clients[dbType]
	if !ok {
		return errors.New("unsupported database")
	}
	c.client = client
	c.dbType = dbType
	c.url = url
	_, err := c.ExecSQL("SELECT 42")
	return err
}

func (c *DB) FetchTables() string {
	var output string
	switch c.dbType {
	case "mysql":
		output, _ = c.ExecSQL("sh")
	case "postgres":
		output, _ = c.ExecSQL("sh")
	case "sqlite":
		output, _ = c.ExecSQL(".schema")
	}
	return output
}

func (c *DB) command(sql string) []string {
	switch c.dbType {
	case "mysql":
		return []string{c.url, "-e", sql}
	case "postgres":
		return []string{c.url, "-c", sql}
	}
	return []string{c.url, sql}
}

func (c *DB) ExecSQL(sql string) (string, error) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	defer signal.Reset(os.Interrupt)

	cmd := exec.Command(c.client, c.command(sql)...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	done := make(chan struct{})
	var output []byte
	var err error
	go func() {
		output, err = cmd.CombinedOutput()
		close(done)
	}()

	// Wait for the command to finish or for a signal to be received
	select {
	case <-done:
	case <-sigChan:
		if runtime.GOOS == "windows" {
			err = cmd.Process.Signal(os.Kill)
		} else {
			err = syscall.Kill(-cmd.Process.Pid, syscall.SIGINT)
		}
	}
	return string(output), err
}
