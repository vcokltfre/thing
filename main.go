package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/vcokltfre/thing/src/commands/cleanup"
	"github.com/vcokltfre/thing/src/commands/pg"
	"github.com/vcokltfre/thing/src/docker"
)

var app = &cli.App{
	Name: "thing",
	Commands: []*cli.Command{
		pg.CommandPg,
		cleanup.CommandCleanup,
	},
}

func main() {
	args := []string{}
	for _, arg := range os.Args {
		if arg == "-c" {
			docker.Client.RemoveAllContainers()

			if len(os.Args) == 2 {
				return
			}

			continue
		}

		args = append(args, arg)
	}

	err := app.Run(args)
	if err != nil {
		fmt.Println(err)
	}
}
