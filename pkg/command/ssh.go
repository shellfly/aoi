package command

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kevinburke/ssh_config"
	"golang.org/x/crypto/ssh"
)

func init() {
	commands["ssh"] = &Ssh{}
}

type Ssh struct {
	client      *ssh.Client
	host        string
	osInfo      string
	isFinished  bool
	initialized bool
}

func (c *Ssh) Name() string {
	return "ssh"
}

func (c *Ssh) Help() string {
	return "/ssh - Generate shell command and execute it on the remote host, e.g. /ssh {host} view listening tcp ports"
}

// Prompt set terminal prompt for ssh command
func (c *Ssh) Prompt(p string) string {
	if c.isFinished {
		return p + ": "
	}
	return fmt.Sprintf("(/ssh %s) %s: ", c.host, p)
}

// Init...
func (c *Ssh) Init(input string) string {
	c.isFinished = true
	index := strings.Index(input, " ")
	if index == -1 {
		if err := c.setHost(input); err == nil {
			fmt.Println("connected to host, you can now ask for command line tasks")
			c.isFinished = false
		}
		return ""
	}

	host, input := input[:index], input[index+1:]
	if err := c.setHost(host); err != nil {
		return ""
	}
	return input
}

func (c *Ssh) IsFinished() bool {
	return c.isFinished
}

func (c *Ssh) Prompts(input string) []string {
	if strings.HasPrefix(input, ":") {
		c.Handle(input[1:])
		return nil
	}

	if c.initialized {
		return []string{input}
	}
	c.initialized = true
	return []string{
		fmt.Sprintf(`
I want you to act as a terminal. I will ask you a question and you will reply with one-line command to do it, avoid pipeline if possible.
I want you to only reply with the code, and nothing else. do not write explanations.
My question is how to %s on %s?
		`, input, c.osInfo),
	}
}

func (c *Ssh) Finish() {
	c.host = ""
	c.initialized = false
	c.isFinished = true
	c.client.Close()
}

// Handle execute command on c.host
func (c *Ssh) Handle(reply string) {
	fmt.Println(reply)
	fmt.Println()
	session, err := c.client.NewSession()
	if err != nil {
		fmt.Println("failed to create session on host: ", err)
		return
	}
	defer session.Close()

	out, err := session.CombinedOutput(reply)
	if err != nil {
		fmt.Println("Failed to run command: ", err)
	}
	fmt.Println(string(out))
	fmt.Println()
}

func (c *Ssh) setHost(host string) error {
	hostname := ssh_config.Get(host, "HostName")
	port := ssh_config.Get(host, "Port")
	user := ssh_config.Get(host, "User")
	identityFiles := ssh_config.GetAll(host, "IdentityFile")
	signer, err := parseSshKey(identityFiles)
	if err != nil {
		fmt.Println("failed to get ssh key: ", err)
		return err
	}
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Minute,
	}
	addr := fmt.Sprintf("%s:%s", hostname, port)
	fmt.Println("connecting to host...")
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		fmt.Println("failed to connect to host: ", err)
		return err
	}

	c.client = client
	c.host = host
	c.osInfo = c.setOSinfo()
	return nil
}

func (c *Ssh) setOSinfo() string {
	session, err := c.client.NewSession()
	if err != nil {
		return ""
	}
	defer session.Close()
	output, err := session.Output("lsb_release -d -s")
	if err == nil {
		return string(output)
	}

	session2, err := c.client.NewSession()
	if err != nil {
		return ""
	}
	defer session2.Close()
	output, err = session2.Output("sw_vers -productVersion")
	if err == nil {
		return fmt.Sprintf("Mac OS %s", output)
	}

	// TODO: support windows
	return ""
}

func parseSshKey(files []string) (ssh.Signer, error) {
	if len(files) == 0 {
		return nil, errors.New("empty file list")
	}
	path := files[0]
	if strings.HasPrefix(path, "~/") {
		dirname, _ := os.UserHomeDir()
		path = filepath.Join(dirname, path[2:])
	}
	key, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ssh.ParsePrivateKey(key)
}
