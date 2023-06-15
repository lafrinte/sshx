package gosh

import (
    "context"
    "fmt"
    "golang.org/x/crypto/ssh"
    "net"
    "strconv"
    "time"
)

type Client struct {
    client  *ssh.Client
    session *ssh.Session
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
        log.Error().Str("action", "ssh").Msgf("failed create ssh client: err -> %s", err)
        return err
    }

    c.client = client

    session, err := client.NewSession()
    if err != nil {
        log.Error().Str("action", "ssh").Msgf("unable to open new session: err -> %s", err)
        return err
    }

    c.session = session

    return nil
}

func (c *Client) Run(ctx context.Context, t *Task, ch chan *TaskState) {
    ctx, cancel := context.WithTimeout(ctx, time.Second*60)
    defer c.Close()
    defer cancel()

    out := make(chan *TaskState, 1)
    go c.run(t, out)

    select {
    case ts := <-out:
        ch <- ts
    case <-ctx.Done():
        ch <- &TaskState{
            ID:     t.ID,
            Bind:   net.JoinHostPort(t.Host, strconv.Itoa(t.Port)),
            Code:   1,
            Stderr: fmt.Sprintf("reach timeout: %ds", 60),
        }
    }
}

func (c *Client) run(t *Task, out chan *TaskState) {
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

    buf, err := c.session.CombinedOutput(t.Command)
    if err != nil {
        ts.Code = 1
        ts.Stderr = fmt.Sprintf("command failed: err -> %s", err)
    }

    ts.Stdout = string(buf)

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
