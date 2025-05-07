package routes

import (
	"chatify/controllers"
	global_types "chatify/types"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func SetupRouters(app *fiber.App) {
	var router fiber.Router = app.Group("/authentication")

	// router.Get("/get-csrf-token", controllers.GetCSRFToken)
	router.Post("/system-login", controllers.SystemLogin)
	// router.Post("/get-refresh-token", controllers.GetRefreshToken)
	// router.Post("/system-logout", middlewares.Authorization(), controllers.SystemLogout)

	app.Use(func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusNotFound).JSON(global_types.IResponseAPI{
			Message:      "This API could not be found",
			ErrorSection: "SetupRouters | search router",
		})
	})
}
