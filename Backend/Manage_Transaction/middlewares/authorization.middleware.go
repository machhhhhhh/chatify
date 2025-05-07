package middlewares

import (
	"chatify/services"
	global_types "chatify/types"
	"chatify/utils"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Authorization() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var token_header string = ctx.Get("Authorization")

		if token_header == "" || strings.Contains(token_header, "Bearer") != true {
			return ctx.Status(http.StatusUnauthorized).JSON(global_types.IResponseAPI{
				Message:      "Unauthorized",
				ErrorSection: "Authorization | validate api header",
			})
		}

		var access_token string = strings.TrimSpace(strings.Replace(token_header, "Bearer", "", 1))

		aes_user, err := utils.AESDecrypted(access_token)
		if err != nil {
			return ctx.Status(http.StatusUnauthorized).JSON(global_types.IResponseAPI{
				Message:      err.Error(),
				ErrorSection: "Authorization | decrypt token",
			})
		}

		account, err := services.FindAccountModel(&aes_user.AccountID)
		if err != nil {
			return ctx.Status(http.StatusUnauthorized).JSON(global_types.IResponseAPI{
				Message:      err.Error(),
				ErrorSection: "GetInformationAccount | find account",
			})
		}

		var body struct {
			Payload string `json:"payload,omitempty"`
		}

		if err := ctx.BodyParser(&body); err != nil {
			if _, ok := err.(*json.UnmarshalTypeError); ok {
				return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
					Message:      "Invalid Payload",
					ErrorSection: "Authorization | check type json",
				})
			}
			return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
				Message:      err.Error(),
				ErrorSection: "Authorization | body parser",
			})
		}

		if body.Payload == "" {
			return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
				Message:      "Empty Payload is not allowed",
				ErrorSection: "Authorization | empty payload",
			})
		}

		payload, err := utils.ReadJWTToken(body.Payload, access_token)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
				Message:      err.Error(),
				ErrorSection: "Authorization | read jwt",
			})
		}

		ctx.Locals("user", account)
		ctx.Locals("access_token", access_token)
		ctx.Locals("payload", payload)
		return ctx.Next()
	}
}
