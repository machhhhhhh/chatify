package services

import (
	"chatify/models"
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
)

func GetContextUser(ctx *fiber.Ctx) (models.Account, error) {
	user := ctx.Locals("user")

	var account models.Account

	data, err := json.Marshal(user)
	if err != nil {
		return account, err
	}

	err = json.Unmarshal(data, &account)
	if err != nil {
		return account, err
	}

	return account, nil
}

func GetContextAccessToken(ctx *fiber.Ctx) (string, error) {
	access_token := ctx.Locals("access_token")

	var token string

	data, err := json.Marshal(access_token)
	if err != nil {
		return token, err
	}

	err = json.Unmarshal(data, &token)
	if err != nil {
		return token, err
	}

	return token, nil
}

func GetContextPayload[T any](ctx *fiber.Ctx) (T, error) {
	var payload T
	request_body := ctx.Locals("payload")

	data, err := json.Marshal(request_body)
	if err != nil {
		return payload, errors.New("Cannot Marshal Body")
	}

	err = json.Unmarshal(data, &payload)
	if err != nil {
		return payload, errors.New("Invalid Payload")
	}

	return payload, nil
}
