package controllers

import (
	"chatify/constants"
	"chatify/databases"
	"chatify/models"
	"chatify/services"
	global_types "chatify/types"
	"chatify/utils"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SystemLogin(ctx *fiber.Ctx) error {
	var body models.Account

	if err := ctx.BodyParser(&body); err != nil {
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
				Message:      "Invalid Payload",
				ErrorSection: "SystemLogin | check type json",
			})
		}
		return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
			Message:      err.Error(),
			ErrorSection: "SystemLogin | body parser",
		})
	}

	if strings.TrimSpace(body.Username) == "" ||
		strings.TrimSpace(body.Password) == "" {
		return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
			Message:      "Incorrect Parameter",
			ErrorSection: "SystemLogin | validate payload",
		})
	}

	var account models.Account

	var query *gorm.DB = databases.DB.
		Scopes(services.DebugMode).
		Where("account_username = ?", strings.TrimSpace(body.Username)).
		Omit("account_identify_number").
		Find(&account)

	if query.Error != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
			Message:      query.Error.Error(),
			ErrorSection: "SystemLogin | find account",
		})
	}

	if query.RowsAffected == 0 {
		return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
			Message:      "Account is not exist",
			ErrorSection: "SystemLogin | check account exist",
		})
	}

	if err := utils.CheckPassword(body.Password, account.Password); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
			Message:      "Account is not exist",
			ErrorSection: "SystemLogin | check account exist",
		})
	}

	access_token, err := utils.AESEncrypted(&global_types.IObjectAES{AccountID: account.AccountID})
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
			Message:      err.Error(),
			ErrorSection: "SystemLogin | encrypted aes",
		})
	}

	var timestamp time.Time = time.Now()

	var logs_authentication models.Logs_Authentication = models.Logs_Authentication{
		AccountID:               account.AccountID,
		LogsActivity:            constants.LogsActivityAuthenticationLogIn,
		ReceiveNotificationAt:   nil,
		LogsNotificationCounter: 0,
		CreatedBy:               account.AccountID,
		CreatedAt:               &timestamp,
		ExpiredAt:               nil,
	}

	//! *******************************************
	//! ********* TRANSACTION DATABASE ************
	//! *******************************************

	var transaction *gorm.DB = databases.DB.Begin()

	if err := transaction.
		Scopes(services.DebugMode).
		Model(&models.Account{}).
		Where("account_id = ?", account.AccountID).
		Update("last_login", time.Now()).
		Error; err != nil {
		transaction.Rollback()
		return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
			Message:      err.Error(),
			ErrorSection: "SystemLogin | update last login",
		})
	}

	if err := transaction.Scopes(services.DebugMode).Create(&logs_authentication).Error; err != nil {
		transaction.Rollback()
		return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
			Message:      err.Error(),
			ErrorSection: "MigrateAccount | create logs authentication",
		})
	}

	transaction.Commit()

	return ctx.Status(fiber.StatusOK).JSON(global_types.IResponseAPI{
		Message: "Login Successfully",
		Data:    access_token,
	})
}
