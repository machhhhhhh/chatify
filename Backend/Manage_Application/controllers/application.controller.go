package controllers

import (
	"chatify/databases"
	"chatify/models"
	"chatify/services"
	global_types "chatify/types"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func MigrateApplication(ctx *fiber.Ctx) error {
	// userMe, err := services.GetContextUser(ctx)
	// if err != nil {
	// 	return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
	// 		Message:      err.Error(),
	// 		ErrorSection: "MigrateApplication | get user from token",
	// 	})
	// }

	// if utils.IsAdmin(userMe.AccountRole) != true {
	// 	return ctx.Status(http.StatusForbidden).JSON(global_types.IResponseAPI{
	// 		Message:      "No Permission",
	// 		ErrorSection: "MigrateApplication | validate role user",
	// 	})
	// }

	// application, err := services.FindApplicationModel()
	// if err != nil {
	// 	return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
	// 		Message:      err.Error(),
	// 		ErrorSection: "MigrateApplication | find application",
	// 	})
	// }

	// if err := services.IsApplicationUpdatingService(&application); err != nil {
	// 	return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
	// 		Message:      err.Error(),
	// 		ErrorSection: "MigrateApplication | check application update schedule",
	// 	})
	// }

	//! *******************************************
	//! ********* TRANSACTION DATABASE ************
	//! *******************************************

	var transaction *gorm.DB = databases.DB.Begin()

	if err := transaction.Scopes(services.DebugMode).AutoMigrate(
		&models.Application{},
	); err != nil {
		transaction.Rollback()
		return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
			Message:      err.Error(),
			ErrorSection: "MigrateApplication | migrate data",
		})
	}

	transaction.Commit()
	transaction = databases.DB.Begin()

	var app models.Application

	query := transaction.Scopes(services.DebugMode).Find(&app)

	if query.Error != nil {
		transaction.Rollback()
		return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
			Message:      query.Error.Error(),
			ErrorSection: "MigrateApplication | find data",
		})
	}

	if query.RowsAffected == 0 {
		var timestamp time.Time = time.Now()
		const start_time string = "01:00"
		const end_time string = "02:00"
		const layout string = "15:04"

		start, err := time.Parse(layout, start_time)
		if err != nil {
			transaction.Rollback()
			return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
				Message:      query.Error.Error(),
				ErrorSection: "MigrateApplication | convert start time",
			})
		}

		end, err := time.Parse(layout, end_time)
		if err != nil {
			transaction.Rollback()
			return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
				Message:      query.Error.Error(),
				ErrorSection: "MigrateApplication | convert end time",
			})
		}

		new_application := models.Application{
			ApplicationName:                    "chatify Management System",
			ApplicationScheduleUpdateStartTime: &start,
			ApplicationScheduleUpdateEndTime:   &end,
			CreatedBy:                          0,
			// CreatedBy:                          userMe.AccountID,
			CreatedAt: &timestamp,
		}

		if err := transaction.Scopes(services.DebugMode).Create(&new_application).Error; err != nil {
			transaction.Rollback()
			return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
				Message:      err.Error(),
				ErrorSection: "MigrateApplication | create application",
			})
		}
	}

	transaction.Commit()

	return ctx.Status(fiber.StatusOK).JSON(global_types.IResponseAPI{
		Message: "Migrate Application Successfully",
	})
}

func GetInformationApplication(ctx *fiber.Ctx) error {
	var application global_types.IApplication

	query := databases.DB.
		Scopes(services.DebugMode).
		Preload("Creator", services.SelectAccount).
		Preload("Updater", services.SelectAccount).
		Scopes(services.SelectApplication).
		Find(&application)

	if query.Error != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
			Message:      query.Error.Error(),
			ErrorSection: "GetInformationApplication | find application",
		})
	}

	if query.RowsAffected == 0 {
		return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
			Message:      "Application configuration is not exist",
			ErrorSection: "GetInformationApplication | check application exist",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(global_types.IResponseAPI{
		Message: "Get Information Application Successfully",
		Data:    application,
	})
}
