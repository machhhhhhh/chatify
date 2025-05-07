package routes

import (
	"chatify/controllers"
	"chatify/middlewares"
	globals_type "chatify/types"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func SetupRouters(app *fiber.App) {
	var router fiber.Router = app.Group("/manage-account")

	router.Post("/migrate-account", controllers.MigrateAccount)
	router.Post("/get-list-account", middlewares.Authorization(), controllers.GetListAccount)
	router.Post("/get-information-account", middlewares.Authorization(), controllers.GetInformationAccount)
	router.Post("/create-account", controllers.CreateAccount)

	app.Use(func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusNotFound).JSON(globals_type.IResponseAPI{
			Message:      "This API could not be found",
			ErrorSection: "SetupRouters | search router",
		})
	})
}
