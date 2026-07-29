package secure

import (
	"errors"

	"github.com/zalando/go-keyring"
)

const SERVICE_KEY = "xiv-otp-service"

func GetValue(user string) (string, error) {
	val, err := keyring.Get(SERVICE_KEY, user)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return "", nil
		}
		return "", err
	}
	return val, nil
}

func SetValue(user string, value string) error {
	err := keyring.Set(SERVICE_KEY, user, value)
	if err != nil {
		return err
	}
	return nil
}

func ClearValue(user string) error {
	err := keyring.Delete(SERVICE_KEY, user)
	if err != nil {
		return err
	}
	return nil
}

func ClearAll() error {
	err := keyring.DeleteAll(SERVICE_KEY)
	if err != nil {
		return err
	}
	return nil
}
