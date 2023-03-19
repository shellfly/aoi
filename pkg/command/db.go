package command

import (
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

type DB struct {
	dummyCommand

	dbType      string
	client      string
	args        []string
	isFinished  bool
	initialized bool
}

var dbTypes = map[string]string{
	"mysql":   "mysql",
	"psql":    "postgres",
	"sqlite3": "sqlite",
	"sqlite":  "sqlite",
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
	parts := strings.Split(input, " ")
	c.client, c.args = parts[0], parts[1:]
	if dbType, ok := dbTypes[c.client]; ok {
		c.dbType = dbType
	} else {
		c.dbType = c.client
	}

	if output, err := c.ExecSQL("SELECT 42"); err != nil {
		fmt.Println("connect to db error: ", output)
		return ""
	}
	c.isFinished = false
	return ""
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
	c.client = ""
	c.args = nil
	c.isFinished = true
}

// Handle execute SQL
func (c *DB) Handle(reply string) {
	sql := reply
	if strings.Contains(reply, "```") {
		sql = extractCode(reply)
		fmt.Println("reply: ", reply)
		fmt.Println("exact code: ", sql)
	}
	output, err := c.ExecSQL(sql)
	if err != nil {
		fmt.Println(output, err)
	} else {
		fmt.Println(output)
	}
	fmt.Println()
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

func (c *DB) execArgs(sql string) []string {
	args := make([]string, len(c.args))
	copy(args, c.args)
	switch c.dbType {
	case "mysql":
		args = append(args, []string{"-e", sql}...)
	case "postgres":
		args = append(args, []string{"-c", sql}...)
	case "sqlite":
		args = append(args, sql)
	}

	return args
}

func (c *DB) ExecSQL(sql string) (string, error) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	defer signal.Reset(os.Interrupt)

	fmt.Println("exec sql: ", c.client, c.execArgs(sql))
	cmd := exec.Command(c.client, c.execArgs(sql)...)
	cmd.Stdin = os.Stdin

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
			err = syscall.Kill(cmd.Process.Pid, syscall.SIGINT)
		}
	}
	return string(output), err
}
