package validator

import (
	"fmt"

	"github.com/pquerna/otp"
)

func ValidateOtpUrl(value string) error {
	_, err := otp.NewKeyFromURL(value)
	if err != nil {
		return fmt.Errorf("invalid OTP URL: %w", err)
	}
	return nil
}
