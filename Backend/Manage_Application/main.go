package main

import (
	"chatify/configs"
	"chatify/constants"
	"chatify/databases"
	"chatify/middlewares"
	"chatify/models"
	"chatify/routes"
	"chatify/services"
	global_types "chatify/types"
	"context"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"gorm.io/gorm"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if _, err := databases.ConnectPostgresWithRetry(ctx, 10*time.Second); err != nil {
		os.Exit(2)
	}

	var success bool = MigrateApplication()
	if success != true {
		os.Exit(2)
	}

	var app *fiber.App = fiber.New(fiber.Config{
		BodyLimit:         5 * 1024 * 1024,                      // 5 MB
		Prefork:           false,                                // ❌ avoid unless doing CPU-bound ops & can handle the complexity
		ProxyHeader:       fiber.HeaderXForwardedFor,            // ✅ required behind reverse proxies like NGINX
		CaseSensitive:     true,                                 // ✅ enforce strict routing if app benefits from it
		StrictRouting:     configs.ENV.IsProductionMode == true, // ✅ easier for dev; set to true in prod for predictability
		ServerHeader:      "Fiber",                              // ✅ branding/info; fine to leave
		AppName:           "chatify",                            // ✅ useful for logging/debugging
		Immutable:         true,                                 // ✅ performance optimization if not modifying request body
		EnablePrintRoutes: configs.ENV.IsProductionMode != true, // ✅ disable in prod for cleaner logs
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			if err == fiber.ErrRequestTimeout {
				return ctx.Status(fiber.StatusGatewayTimeout).JSON(global_types.IResponseAPI{
					Message:      "Gateway Timeout",
					ErrorSection: "SetupRouters | gateway timeout",
				})
			}
			return fiber.DefaultErrorHandler(ctx, err) // default behavior
		},
	})

	app.Use(timeout.NewWithContext(func(ctx *fiber.Ctx) error {
		// inside this handler pull the context via c.UserContext() and honor it in long-running ops
		return ctx.Next()
	}, 10*time.Second)) // apply a 10s timeout to **all** routes

	app.Use(compress.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowHeaders: strings.Join([]string{"Content-Type", "Content-Length", "Authorization"}, ","),
		AllowOriginsFunc: func(origin string) bool {
			if origin == configs.ENV.ServerSetting.BaseURL {
				return true
			}
			var allow_domain *regexp.Regexp = regexp.MustCompile(`^(?:http|https)://[a-zA-Z][a-zA-Z0-9\-.]+(?::[0-9]{1,5})?$`)
			return allow_domain.MatchString(origin)
		},
		AllowMethods:     strings.Join(constants.GetAllHTTPMethods(), ","),
		AllowCredentials: true,
	}))

	app.Use(middlewares.RateLimiter(100, time.Minute)) // 1 minute per 100 requests for 1 IP

	routes.SetupRouters(app)

	go func() {
		if err := app.Listen(":" + configs.ENV.ServerSetting.Port); err != nil {
			log.Println("🧨 Cannot Start Server:", err)
			os.Exit(1)
		}
	}()

	log.Println("🚀 Server started successfully, running on port:", configs.ENV.ServerSetting.Port)

	<-ctx.Done()
	log.Println("🛑 Gracefully shutting down...")
	if err := app.Shutdown(); err != nil {
		os.Exit(0)
	}
	log.Println("✅ Shutdown complete")
}

func MigrateApplication() bool {
	var transaction *gorm.DB = databases.DB.Begin()

	if err := transaction.Scopes(services.DebugMode).AutoMigrate(
		&models.Application{},
	); err != nil {
		transaction.Rollback()
		return false
	}

	transaction.Commit()
	transaction = databases.DB.Begin()

	var app models.Application

	var query *gorm.DB = transaction.Scopes(services.DebugMode).Find(&app)

	if query.Error != nil {
		transaction.Rollback()
		return false
	}

	if query.RowsAffected == 0 {
		var timestamp time.Time = time.Now()
		const start_time string = "01:00"
		const end_time string = "02:00"
		const layout string = "15:04"

		start, err := time.Parse(layout, start_time)
		if err != nil {
			transaction.Rollback()
			return false
		}

		end, err := time.Parse(layout, end_time)
		if err != nil {
			transaction.Rollback()
			return false
		}

		new_application := models.Application{
			ApplicationName:                    "chatify Management System",
			ApplicationScheduleIsActive:        true,
			ApplicationScheduleUpdateStartTime: &start,
			ApplicationScheduleUpdateEndTime:   &end,
			CreatedBy:                          0,
			CreatedAt:                          &timestamp,
		}

		if err := transaction.Scopes(services.DebugMode).Create(&new_application).Error; err != nil {
			transaction.Rollback()
			return false
		}
	}

	transaction.Commit()

	return true
}
