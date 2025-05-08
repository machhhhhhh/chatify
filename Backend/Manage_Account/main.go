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
	"chatify/utils"
	"context"
	"fmt"
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

	var success bool = MigrateAccount()
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

func MigrateAccount() bool {
	var transaction *gorm.DB = databases.DB.Begin()

	if err := transaction.Scopes(services.DebugMode).AutoMigrate(
		&models.Account{},
		&models.Logs_Authentication{},
	); err != nil {
		transaction.Rollback()
		return false
	}

	transaction.Commit()
	transaction = databases.DB.Begin()

	var account []models.Account

	if err := transaction.Scopes(services.DebugMode).Find(&account).Error; err != nil {
		transaction.Rollback()
		return false
	}

	if len(account) == 0 {
		var timestamp time.Time = time.Now()

		password, err := utils.HashPassword("555")
		if err != nil {
			transaction.Rollback()
			return false
		}

		account = append(account, models.Account{
			AccountNumber:         "AC000001",
			FirstName:             "Super",
			LastName:              "Admin",
			AccountRole:           constants.AccountRoleSuperAdmin,
			IdentifyNumber:        "5555555555555",
			Username:              "super_admin@chatify.com",
			Password:              password,
			Email:                 "super_admin@chatify.com",
			PhoneNumber:           "66830313097",
			IsVerifiedEmail:       true,
			IsVerifiedPhoneNumber: true,
			IsActive:              true,
			CreatedAt:             &timestamp,
		})

		account = append(account, models.Account{
			AccountNumber:         "AC000002",
			FirstName:             "Admin",
			LastName:              "chatify",
			IdentifyNumber:        "5555555555555",
			AccountRole:           constants.AccountRoleAdmin,
			Username:              "admin@chatify.com",
			Password:              password,
			Email:                 "admin@chatify.com",
			PhoneNumber:           "66830313097",
			IsVerifiedEmail:       true,
			IsVerifiedPhoneNumber: true,
			IsActive:              true,
			CreatedAt:             &timestamp,
		})

		for i := range account {
			if err := transaction.Scopes(services.DebugMode).Create(&account[i]).Error; err != nil {
				transaction.Rollback()
				return false
			}

			if err := transaction.
				Scopes(services.DebugMode).
				Model(&models.Account{}).
				Where("account_id = ?", account[i].AccountID).
				Updates(&models.Account{
					AccountNumber: fmt.Sprintf("AC%06d", account[i].AccountID),
					CreatedBy:     account[i].AccountID,
				}).
				Error; err != nil {
				transaction.Rollback()
				return false
			}
		}
	}

	transaction.Commit()

	return true
}
