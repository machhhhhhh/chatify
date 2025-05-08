package routes

import (
	"chatify/controllers"
	"chatify/middlewares"
	globals_type "chatify/types"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func SetupRouters(app *fiber.App) {
	var router fiber.Router = app.Group("/manage-transaction")

	router.Post("/migrate-transaction", controllers.MigrateTransaction)
	router.Post("/get-list-transaction", middlewares.Authorization(), controllers.GetListTransaction)
	router.Post("/get-information-transaction", middlewares.Authorization(), controllers.GetInformationTransaction)
	router.Post("/create-transaction", middlewares.UploadFile(), middlewares.Authorization(), controllers.CreateTransaction)

	app.Use(func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusNotFound).JSON(globals_type.IResponseAPI{
			Message:      "This API could not be found",
			ErrorSection: "SetupRouters | search router",
		})
	})
}
