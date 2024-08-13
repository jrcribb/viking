package machine

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/d3witt/viking/cli/command"
	"github.com/d3witt/viking/config"
	"github.com/urfave/cli/v2"
)

func NewAddCmd(vikingCli *command.Cli) *cli.Command {
	return &cli.Command{
		Name:        "add",
		Usage:       "Add a new machine",
		Description: "This command adds a new machine to the list of machines. No action is taken on the machine itself. Ensure your computer has SSH access to this machine.",
		Args:        true,
		ArgsUsage:   "HOST",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "Machine name",
			},
			&cli.StringFlag{
				Name:    "user",
				Aliases: []string{"u"},
				Value:   "root",
				Usage:   "SSH user name",
			},
			&cli.StringFlag{
				Name:    "key",
				Aliases: []string{"k"},
				Usage:   "SSH key name",
			},
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   22,
			},
		},
		Action: func(ctx *cli.Context) error {
			host := ctx.Args().First()
			name := ctx.String("name")
			user := ctx.String("user")
			key := ctx.String("key")
			port := ctx.Int("port")

			return runAdd(vikingCli, host, port, name, user, key)
		},
	}
}

func runAdd(vikingCli *command.Cli, host string, port int, name, user, key string) error {
	if name == "" {
		name = command.GenerateRandomName()
	}

	if key != "" {
		_, err := vikingCli.Config.GetKeyByName(key)
		if err != nil {
			return err
		}
	}

	if h, p, err := net.SplitHostPort(host); err == nil {
		host = h
		port, err = strconv.Atoi(p)
		if err != nil {
			return err
		}
	}

	hostIp := net.ParseIP(host)
	if hostIp == nil {
		return errors.New("host must be valid ip address")
	}

	if err := vikingCli.Config.AddMachine(config.Machine{
		Name:      name,
		Host:      []net.IP{hostIp},
		Port:      port,
		User:      user,
		Key:       key,
		CreatedAt: time.Now(),
	}); err != nil {
		return err
	}

	fmt.Fprintf(vikingCli.Out, "Machine %s added.\n", name)

	return nil
}
