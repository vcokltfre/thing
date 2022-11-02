package pg

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/urfave/cli/v2"
	"github.com/vcokltfre/thing/src/docker"
)

var CommandPg = &cli.Command{
	Name:  "pg",
	Usage: "Start PostgreSQL containers",
	Action: func(c *cli.Context) error {
		rand.Seed(time.Now().UnixNano())

		port, _ := nat.NewPort("tcp", "5432")
		hostPort := nat.PortBinding{
			HostIP:   "0.0.0.0",
			HostPort: fmt.Sprintf("%d", rand.Intn(10000)+10000),
		}

		cid, err := docker.Client.Start(
			"docker.io/library/postgres:latest",
			"postgres:latest",
			fmt.Sprintf("pg__%d", rand.Int()),
			[]string{
				"POSTGRES_PASSWORD=postgres",
				"POSTGRES_USER=postgres",
				fmt.Sprintf("POSTGRES_DB=%s", c.Args().First()),
			},
			nat.PortMap{
				port: []nat.PortBinding{hostPort},
			},
		)
		if err != nil {
			return err
		}

		address := fmt.Sprintf("postgresql://postgres:postgres@localhost:%s/%s", hostPort.HostPort, c.Args().First())

		fmt.Printf("Started container %s with address %s\n\n\u001b[0mPress ENTER to stop.\n", cid[:7], address)
		fmt.Scanln()

		err = docker.Client.StopContainer(cid)
		if err != nil {
			return err
		}

		return docker.Client.RemoveContainer(cid)
	},
}
