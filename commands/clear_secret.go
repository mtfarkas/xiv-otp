package commands

import (
	"context"
	"fmt"

	"github.com/mfarkas/xiv-otp/secure"
	"github.com/mfarkas/xiv-otp/validator"
	"github.com/urfave/cli/v3"
)

func NewClearSecretCommand() *cli.Command {
	return &cli.Command{
		Name:        "clear",
		Description: "Clear OTP Auth secret value for a single account",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:      "account",
				Required:  false,
				Value:     "default",
				Aliases:   []string{"a"},
				Validator: validator.NonEmptyString,
				Usage:     "The name of the user account to manage the secret for",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			account := c.String("account")

			err := secure.ClearValue(account)
			if err != nil {
				return fmt.Errorf("error clearing '%s' from secret store: %w", account, err)
			}
			return nil
		},

		Commands: []*cli.Command{
			{
				Name:        "all",
				Description: "Clear OTP Auth secret values for all accounts",
				Action: func(ctx context.Context, c *cli.Command) error {
					err := secure.ClearAll()
					if err != nil {
						return fmt.Errorf("error clearing secret store: %w", err)
					}
					return nil
				},
			},
		},
	}
}
