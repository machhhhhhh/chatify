package utils

import (
	"errors"
	"strings"

	"github.com/matthewhartstonge/argon2"
)

func HashPassword(password string) (string, error) {
	var argon argon2.Config = argon2.DefaultConfig()
	encoded, err := argon.HashEncoded([]byte(strings.TrimSpace(password)))
	return string(encoded), err
}

func CheckPassword(input_password string, compare_password string) error {
	is_match, err := argon2.VerifyEncoded([]byte(strings.TrimSpace(input_password)), []byte(strings.TrimSpace(compare_password)))
	if err != nil {
		return err
	}

	if is_match != true {
		return errors.New("Not match!")
	}

	return nil
}
