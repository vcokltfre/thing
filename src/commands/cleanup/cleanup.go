package cleanup

import (
	"github.com/urfave/cli/v2"
	"github.com/vcokltfre/thing/src/docker"
)

var CommandCleanup = &cli.Command{
	Name:  "cleanup",
	Usage: "Remove all containers and images",
	Action: func(c *cli.Context) error {
		return docker.Client.RemoveAllContainers()
	},
}
