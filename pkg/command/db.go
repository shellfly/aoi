package command

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rest-go/rest/pkg/sql"
)

func init() {
	commands["db"] = &DB{}
}

type DB struct {
	dummyCommand

	dbType      string
	db          *sql.DB
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
		if err := c.setDB(input); err == nil {
			fmt.Println("connected to db, you can now ask for SQL tasks")
			c.isFinished = false
		} else {
			fmt.Println(err)
			fmt.Println()
		}
		return ""
	}

	host, input := input[:index], input[index+1:]
	if err := c.setDB(host); err != nil {
		fmt.Println(err)
		fmt.Println()
		return ""
	}
	return input
}

func (c *DB) IsFinished() bool {
	return c.isFinished
}

func (c *DB) Prompts(input string) []string {
	if strings.HasPrefix(input, ":") {
		c.Handle(input[1:])
		return nil
	}

	if c.initialized {
		return []string{input}
	}
	c.initialized = true

	prompts := []string{"Given these table definitions: \n"}
	tables := c.db.FetchTables()
	definitions := make([]string, 0, len(tables))
	for _, table := range tables {
		definitions = append(definitions, table.String())
	}
	prompts = append(prompts, fmt.Sprintf("%s\n", strings.Join(definitions, "\n\n")))
	prompts = append(prompts, fmt.Sprintf("You are a %s expert, give SQL for %s , reply with the code, and nothing else.", c.dbType, input))
	return prompts
}

func (c *DB) Finish() {
	c.dbType = ""
	c.db = nil
	c.isFinished = true
}

// Handle execute SQL
func (c *DB) Handle(reply string) {
	sql := reply
	if strings.Contains(reply, "```") {
		sql = extractCode(reply)
	}
	fmt.Println(sql)
	fmt.Println()
	output, err := c.ExecSQL(sql)
	if err != nil {
		fmt.Println(output, err)
	} else {
		fmt.Println(output)
	}
	fmt.Println()
}

func (c *DB) setDB(url string) error {
	db, err := sql.Open(url)
	if err != nil {
		return err
	}
	c.db = db
	c.dbType = db.DriverName
	return nil
}

func (c *DB) ExecSQL(sql string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	data, err := c.db.FetchData(ctx, sql)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", data), nil
}
