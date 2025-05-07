package utils

import (
	"strings"

	"github.com/matthewhartstonge/argon2"
)

func HashPassword(password string) (string, error) {
	var argon argon2.Config = argon2.DefaultConfig()
	encoded, err := argon.HashEncoded([]byte(strings.TrimSpace(password)))
	return string(encoded), err
}

func CheckPassword(input_password string, compare_password string) (bool, error) {
	return argon2.VerifyEncoded([]byte(strings.TrimSpace(input_password)), []byte(compare_password))
}
