package pool

import (
	"context"
	"fmt"
	"github.com/bytedance/gopkg/util/gopool"
	"github.com/bytedance/sonic"
	"net"
	"sshx/gosh/ext/b64"
	"strconv"
	"time"
)

var (
	Pool gopool.Pool
)

func init() {
	Pool = gopool.NewPool("trigger", 1024, gopool.NewConfig())
}

type TaskTrigger struct {
	raw   string
	t     []*Task
	ctx   context.Context
	in    chan *Task
	out   chan *TaskState
	state TasksState
}

type Task struct {
	ID       string `json:"id"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	UserName string `json:"username"`
	Password string `json:"password"`
	Command  string `json:"command"`
	Handler  string `json:"handler"`
}

type TasksState []*TaskState

type TaskState struct {
	ID     string `json:"id"`
	Bind   string `json:"bind"`
	Stderr string `json:"stderr"`
	Stdout string `json:"stdout"`
	Code   int    `json:"code"`
}

/*
Unmarshal use to decode the task string to struct
*/
func (t *TaskTrigger) Unmarshal() error {
	buf, err := b64.Decrypt(t.raw)
	if err != nil {
		return fmt.Errorf("failed decrypt by base64, err -> %s", err)
	}

	err = sonic.Unmarshal(buf, &t.t)
	if err != nil {
		return fmt.Errorf("failed unmarshal by json, err -> %s", err)
	}

	t.in = make(chan *Task, len(t.t))
	t.out = make(chan *TaskState, len(t.t))

	// charge task into channel
	for _, task := range t.t {
		t.in <- task
	}

	return nil
}

// display used to display the task state data in json in terminal
func (t *TaskTrigger) display() error {
	raw, err := sonic.Marshal(t.state)
	if err != nil {
		return fmt.Errorf("failed unmarshal by json, err -> %s", err)
	}

	fmt.Println(string(raw))

	return nil
}

func (t *TaskTrigger) run() {
	var (
		maxLen     = len(t.in)
		maxWaitSec = 10
	)

	for i := 1; i <= maxLen; i++ {
		Pool.CtxGo(t.ctx, func() {
			tk := <-t.in

			ctx, cancel := context.WithTimeout(t.ctx, time.Duration(maxWaitSec)*time.Second)
			defer cancel()

			out := make(chan *TaskState, 1)
			Pool.CtxGo(t.ctx, func() {
				client := &Client{}
				client.Run(tk, out)
			})

			select {
			case ret := <-out:
				t.out <- ret
				return
			case <-ctx.Done():
				t.out <- &TaskState{
					ID:     tk.ID,
					Bind:   net.JoinHostPort(tk.Host, strconv.Itoa(tk.Port)),
					Code:   1,
					Stderr: fmt.Sprintf("reach timeout: %ds", maxWaitSec),
				}
				return
			}
		})
	}
}

func (t *TaskTrigger) Run() error {
	var (
		maxWaitSec = 30
		length     = len(t.in)
	)

	t.run()

	ctx, cancel := context.WithTimeout(t.ctx, time.Duration(maxWaitSec)*time.Second)
	defer cancel()

	for {
		if len(t.state) == length {
			break
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("reach timeout for tasks: %ds. any task may hang and not exit yet", maxWaitSec)
		case tk := <-t.out:
			t.state = append(t.state, tk)
		}
	}

	return t.display()
}

func New(j string) (*TaskTrigger, error) {
	tj := &TaskTrigger{
		raw: j,
		ctx: context.Background(),
	}

	if err := tj.Unmarshal(); err != nil {
		return nil, err
	}

	return tj, nil
}
