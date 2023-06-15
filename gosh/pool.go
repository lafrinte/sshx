package gosh

import (
    "context"
    "fmt"
    "github.com/bytedance/gopkg/util/gopool"
    "github.com/bytedance/sonic"
    "time"
)

var Pool gopool.Pool

func init() {
    Pool = gopool.NewPool("trigger", 1024, gopool.NewConfig())
}

type TaskTrigger struct {
    task string
    ctx  context.Context
    out  chan *TaskState
}

type Tasks struct {
    ID    string  `json:"id"`
    Tasks []*Task `json:"tasks"`
}

type Task struct {
    ID       string `json:"id"`
    Host     string `json:"host"`
    Port     int    `json:"port"`
    UserName string `json:"username"`
    Password string `json:"password"`
    Command  string `json:"command"`
}

type TasksState struct {
    ID     string       `json:"id"`
    States []*TaskState `json:"states"`
}

type TaskState struct {
    ID     string `json:"id"`
    Bind   string `json:"bind"`
    Stderr string `json:"stderr"`
    Stdout string `json:"stdout"`
    Code   int    `json:"code"`
}

func (t *TaskTrigger) Run() error {
    var (
        out = &TasksState{}
    )

    buf, err := B64Decrypt(t.task)
    if err != nil {
        log.Error().Str("action", "base64").Msgf("failed decrypt by base64: err -> %s", err)
        return fmt.Errorf("failed decrypt by base64")
    }

    tasks := &Tasks{}
    err = sonic.Unmarshal(buf, tasks)
    if err != nil {
        log.Error().Str("action", "json").Msgf("failed unmarshal by json: err -> %s", err)
        return fmt.Errorf("failed unmarshal by json")
    }

    out.ID = tasks.ID

    ctx, cancel := context.WithCancel(t.ctx)
    defer cancel()

    for _, task := range tasks.Tasks {
        Pool.CtxGo(ctx, func() {
            client := &Client{}
            client.run(task, t.out)
        })
    }

LOOP:
    for {
        if len(out.States) == len(tasks.Tasks) {
            break
        }

        select {
        case <-time.After(time.Second * 600):
            break LOOP
        case t := <-t.out:
            out.States = append(out.States, t)
        }
    }

    raw, err := sonic.Marshal(out)
    if err != nil {
        log.Error().Str("action", "json").Msgf("failed marshal by json: err -> %s", err)
        return err
    }

    fmt.Println(string(raw))
    return nil
}

func NewTaskTrigger(j string) *TaskTrigger {
    return &TaskTrigger{
        j,
        context.Background(),
        make(chan *TaskState),
    }
}
