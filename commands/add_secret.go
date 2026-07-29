package commands

import (
	"context"
	"fmt"

	"github.com/mfarkas/xiv-otp/secure"
	"github.com/mfarkas/xiv-otp/validator"
	"github.com/urfave/cli/v3"
)

func NewAddSecretCommand() *cli.Command {
	return &cli.Command{
		Name:        "add",
		Description: "Add OTP Auth secret values",
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
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:      "secret",
				UsageText: "\"otpauth://totp/...\"",
				Config: cli.StringConfig{
					TrimSpace: true,
				},
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			secretValue := c.StringArg("secret")
			accountValue := c.String("account")

			err := validator.ValidateOtpUrl(secretValue)
			if err != nil {
				return fmt.Errorf("invalid OTP URL format: %w", err)
			}

			err = secure.SetValue(accountValue, secretValue)
			if err != nil {
				return fmt.Errorf("error setting secret value: %w", err)
			}
			return nil
		},
	}
}
