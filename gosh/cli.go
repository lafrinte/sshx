package gosh

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/urfave/cli/v2"
	"sshx/gosh/ext/b64"
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

func SSHTest() *cli.Command {
	return &cli.Command{
		Name:  "test",
		Usage: "test ssh protocol connection.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "host",
				Usage:       "ssh host",
				Value:       "localhost",
				DefaultText: "localhost",
			},
			&cli.StringFlag{
				Name:    "user",
				Aliases: []string{"u"},
				Usage:   "ssh user",
				Value:   fs.UserName(),
			},
			&cli.StringFlag{
				Name:     "password",
				Aliases:  []string{"p"},
				Usage:    "ssh password",
				Required: true,
			},
			&cli.IntFlag{
				Name:        "port",
				Aliases:     []string{"P"},
				Usage:       "ssh port",
				Value:       22,
				DefaultText: "22",
			},
			&cli.StringFlag{
				Name:     "command",
				Aliases:  []string{"c"},
				Usage:    "command",
				Required: true,
			},
			&cli.StringFlag{
				Name:        "type",
				Aliases:     []string{"t"},
				Usage:       "interrupter type. sh | bash | py2 | py3 | ruby",
				Value:       "bash",
				DefaultText: "bash",
			},
		},
		Action: func(context *cli.Context) error {
			t := []*pool.Task{
				{
					ID:       "test",
					Host:     context.String("host"),
					Port:     context.Int("port"),
					UserName: context.String("user"),
					Password: context.String("password"),
					Handler:  context.String("type"),
					Command:  context.String("command"),
				},
			}

			decode, err := sonic.Marshal(t)
			if err != nil {
				return cli.Exit(fmt.Sprintf("json marshal error: err -> %s", err), 2)
			}

			ciphers := b64.Encrypt(decode)
			trigger, err := pool.New(ciphers)
			if err != nil {
				return cli.Exit(fmt.Sprintf("init error: err -> %s", err), 2)
			}

			if err := trigger.Run(); err != nil {
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
