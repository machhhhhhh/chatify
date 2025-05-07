package routes

import (
	"chatify/controllers"
	global_types "chatify/types"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func SetupRouters(app *fiber.App) {
	var router fiber.Router = app.Group("/authentication")

	router.Post("/system-login", controllers.SystemLogin)

	app.Use(func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusNotFound).JSON(global_types.IResponseAPI{
			Message:      "This API could not be found",
			ErrorSection: "SetupRouters | search router",
		})
	})
}
