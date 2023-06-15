package main

import (
    "github.com/urfave/cli/v2"
    "os"
    "sshx/gosh"
)

func InitCLI() *cli.App {
    app := cli.NewApp()
    app.UseShortOptionHandling = true

    app.Name = "gosh"
    app.Version = "v1.0.0-alpha"
    app.Commands = []*cli.Command{
        gosh.ParallelSSH(),
        gosh.Version(),
    }

    return app
}

func main() {
    app := InitCLI()
    if err := app.Run(os.Args); err != nil {
        os.Exit(1)
    }

    os.Exit(0)
}
