package gosh

import (
    "fmt"
    "github.com/urfave/cli/v2"
)

func ParallelSSH() *cli.Command {
    return &cli.Command{
        Name:  "pssh",
        Usage: "parallel ssh protocol connection.",
        Flags: []cli.Flag{
            &cli.StringFlag{
                Name:     "json",
                Required: true,
            },
        },
        Action: func(context *cli.Context) error {
            t := NewTaskTrigger(context.String("json"))
            if err := t.Run(); err != nil {
                return cli.Exit("", 1)
            }
            return nil
        },
    }
}

func Version() *cli.Command {
    return &cli.Command{
        Name: "version",
        Action: func(c *cli.Context) error {
            fmt.Println("version: v1.0.0-alpha")
            return nil
        },
        Usage: "display version information.",
    }
}
