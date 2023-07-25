package gosh

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"sshx/gosh/ext/pool"
	"sshx/gosh/fs"
)

func ParallelSSH() *cli.Command {
	return &cli.Command{
		Name:  "pssh",
		Usage: "parallel ssh protocol connection.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "json",
				Aliases: []string{"j"},
				Usage:   "base64 encoded json string",
			},
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Usage:   "file path which contains base64 encoded json string",
			},
		},
		Action: func(context *cli.Context) error {
			var (
				f   = context.String("file")
				j   = context.String("json")
				raw string
			)

			if f == "" && j == "" {
				return cli.Exit("argument required. please pass -f or -j", 1)
			}

			if j != "" {
				raw = j
			}

			if f != "" {
				s, err := fs.ReadFile(f)
				if err != nil {
					return cli.Exit(fmt.Sprintf("io error occured while reading file: path -> %s", f), 1)
				}

				raw = s
			}

			t, err := pool.New(raw)
			if err != nil {
				return cli.Exit(fmt.Sprintf("init error: err -> %s", err), 2)
			}
			if err := t.Run(); err != nil {
				return cli.Exit(fmt.Sprintf("timeout or unmarshal failed: err -> %s", err), 1)
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
