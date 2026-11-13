package commands

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mtfarkas/xiv-otp/secure"
	"github.com/mtfarkas/xiv-otp/validator"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/urfave/cli/v3"
)

func NewAuthCommand() *cli.Command {
	return &cli.Command{
		Name:        "auth",
		Description: "Perform the login flow for a specific account, or with a specified OTP Auth URL",
		Flags: []cli.Flag{
			&cli.Uint16Flag{
				Name:  "port",
				Value: 4646,
			},
		},
		Commands: []*cli.Command{
			{
				Name:        "login",
				Description: "Perform the login flow for a specific account",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:      "account",
						Value:     "default",
						Aliases:   []string{"a"},
						Validator: validator.NonEmptyString,
						Usage:     "The name of the user account to use the secret for",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					account := c.String("account")
					port := c.Uint16("port")

					otpUrl, err := secure.GetValue(account)
					if err != nil {
						return fmt.Errorf("login failed; failed to get auth secret from keyring: %w", err)
					}

					return performAuth(otpUrl, port)
				},
			},
			{
				Name:        "login-url",
				Description: "Perform the login flow with a specified OTP Auth URL",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:      "url",
						UsageText: "\"otpauth://totp/...\"",
						Config: cli.StringConfig{
							TrimSpace: true,
						},
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					otpUrl := c.StringArg("url")
					port := c.Uint16("port")

					err := validator.ValidateOtpUrl(otpUrl)
					if err != nil {
						return fmt.Errorf("invalid OTP URL format: %w", err)
					}

					return performAuth(otpUrl, port)
				},
			},
		},
	}
}

func performAuth(otpUrl string, port uint16) error {
	if len(otpUrl) < 1 {
		return fmt.Errorf("failed to log in; OTP Auth URL was empty")
	}

	key, err := otp.NewKeyFromURL(otpUrl)
	if err != nil {
		return fmt.Errorf("failed to get key from URL: %w", err)
	}

	const maxRetries = 3
	var lastError error

	for attempt := range maxRetries {
		currentOtp, err := totp.GenerateCode(key.Secret(), time.Now())
		if err != nil {
			return fmt.Errorf("failed to generate TOTP key from secret: %w", err)
		}

		requestUrl := fmt.Sprintf("http://localhost:%d/ffxivlauncher/%s", port, currentOtp)

		_, lastError = http.Get(requestUrl)
		if lastError == nil {
			return nil
		}

		if attempt < maxRetries-1 {
			time.Sleep(2 * time.Second)
		}
	}

	if lastError != nil {
		fmt.Println("Couldn't send TOTP key to XIVLauncher. Ensure it's running, it's waiting for a TOTP for the login, and OTP macro support is turned on.")
		fmt.Println()

		currentOtp, err := totp.GenerateCode(key.Secret(), time.Now())
		if err != nil {
			fmt.Println("No fallback TOTP available.")
		} else {
			fmt.Printf("Fallback TOTP code to paste: %s", currentOtp)
			fmt.Println()
		}

		return fmt.Errorf("failed to log in: %w", lastError)
	}

	return nil
}
