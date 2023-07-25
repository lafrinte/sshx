package pool

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Client struct {
	client  *ssh.Client
	session *ssh.Session
}

type singleWriter struct {
	b  bytes.Buffer
	mu sync.Mutex
}

func (w *singleWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.b.Write(p)
}

func (c *Client) DialWithPassword(host string, port int, username string, password string) error {
	client, err := ssh.Dial("tcp", net.JoinHostPort(host, strconv.Itoa(port)), &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
		Config: ssh.Config{
			Ciphers: []string{
				"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-cbc", "aes192-cbc",
				"aes256-cbc", "3des-cbc", "des-cbc", "arcfour128", "arcfour256",
			},
		},
	})

	if err != nil {
		return err
	}

	c.client = client

	session, err := client.NewSession()
	if err != nil {
		return err
	}

	c.session = session

	return nil
}

func (c *Client) Run(t *Task, out chan *TaskState) {
	ts := &TaskState{
		ID:   t.ID,
		Bind: net.JoinHostPort(t.Host, strconv.Itoa(t.Port)),
	}

	if c.client == nil || c.session == nil {
		if err := c.DialWithPassword(t.Host, t.Port, t.UserName, t.Password); err != nil {
			ts.Code = 1
			ts.Stderr = fmt.Sprintf("failed connection: err -> %s", err)
			out <- ts
			return
		}
	}

	switch OperatorEnum(t.Handler) {
	case Python, Python3:
		var b singleWriter

		c.session.Stdout = &b
		c.session.Stderr = &b
		c.session.Stdin = strings.NewReader(strings.ReplaceAll(t.Command, "\t", "    "))
		err := c.session.Run(OperatorEnum(t.Handler))
		if err != nil {
			ts.Code = 1
			ts.Stderr = fmt.Sprintf("command failed: err -> %s", err)
		}

		ts.Stdout = string(b.b.Bytes())
	default:
		buf, err := c.session.CombinedOutput(t.Command)
		if err != nil {
			ts.Code = 1
			ts.Stderr = fmt.Sprintf("command failed: err -> %s", err)
		}

		ts.Stdout = string(buf)
	}

	out <- ts
}

func (c *Client) Close() error {
	if err := c.session.Close(); err != nil {
		return err
	}

	if err := c.Close(); err != nil {
		return err
	}

	return nil
}
