package validator

import (
	"fmt"
	"strings"
)

func NonEmptyString(s string) error {
	trimmed := strings.TrimSpace(s)
	if len(trimmed) < 1 {
		return fmt.Errorf("string value cannot be empty")
	}
	return nil
}
