package utils

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWTToken(payload any, secret string) (string, error) {
	var claims *jwt.Token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"payload": payload,
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	})

	token, err := claims.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func ReadJWTToken(payload string, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(strings.TrimSpace(payload), func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	if err != nil || token.Valid != true {
		return nil, errors.New("Invalid Token")
	}

	// global value result from client
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok != true {
		return claims, errors.New("Cannot claim the payload token")
	}

	return claims, nil
}
