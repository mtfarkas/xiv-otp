package main

import (
	"context"
	"log"
	"os"

	"github.com/mfarkas/xiv-otp/commands"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := cli.Command{
		Commands: []*cli.Command{
			{
				Name:        "secret",
				Description: "Manage OTP Auth secrets",
				Commands: []*cli.Command{
					commands.NewAddSecretCommand(),
					commands.NewClearSecretCommand(),
				},
			},
			commands.NewAuthCommand(),
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
