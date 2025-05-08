package controllers

import (
	"chatify/constants"
	"chatify/databases"
	"chatify/models"
	"chatify/services"
	global_types "chatify/types"
	"chatify/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func MigrateAccount(ctx *fiber.Ctx) error {
	// userMe, err := services.GetContextUser(ctx)
	// if err != nil {
	// 	return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
	// 		Message:      err.Error(),
	// 		ErrorSection: "MigrateAccount | get user from token",
	// 	})
	// }

	// if utils.IsAdmin(userMe.AccountRole) != true {
	// 	return ctx.Status(http.StatusForbidden).JSON(global_types.IResponseAPI{
	// 		Message:      "No Permission",
	// 		ErrorSection: "MigrateAccount | validate role user",
	// 	})
	// }

	// application, err := services.FindApplicationModel()
	// if err != nil {
	// 	return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
	// 		Message:      err.Error(),
	// 		ErrorSection: "MigrateAccount | find application",
	// 	})
	// }

	// if err := services.IsApplicationUpdatingService(&application); err != nil {
	// 	return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
	// 		Message:      err.Error(),
	// 		ErrorSection: "MigrateAccount | check application update schedule",
	// 	})
	// }

	//! *******************************************
	//! ********* TRANSACTION DATABASE ************
	//! *******************************************

	var transaction *gorm.DB = databases.DB.Begin()

	if err := transaction.Scopes(services.DebugMode).AutoMigrate(
		&models.Account{},
		&models.Logs_Authentication{},
	); err != nil {
		transaction.Rollback()
		return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
			Message:      err.Error(),
			ErrorSection: "MigrateAccount | migrate data",
		})
	}

	transaction.Commit()
	transaction = databases.DB.Begin()

	var account []models.Account

	if err := transaction.Scopes(services.DebugMode).Find(&account).Error; err != nil {
		transaction.Rollback()
		return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
			Message:      err.Error(),
			ErrorSection: "MigrateAccount | find account",
		})
	}

	if len(account) == 0 {
		var timestamp time.Time = time.Now()

		password, err := utils.HashPassword("555")
		if err != nil {
			transaction.Rollback()
			return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
				Message:      err.Error(),
				ErrorSection: "MigrateAccount | hashing password",
			})
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
				return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
					Message:      err.Error(),
					ErrorSection: "MigrateAccount | create account",
				})
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
				return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
					Message:      err.Error(),
					ErrorSection: "MigrateAccount | update account number",
				})
			}
		}
	}

	transaction.Commit()

	return ctx.Status(fiber.StatusOK).JSON(global_types.IResponseAPI{
		Message: "Migrate Account Successfully",
	})
}

func GetListAccount(ctx *fiber.Ctx) error {
	userMe, err := services.GetContextUser(ctx)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
			Message:      err.Error(),
			ErrorSection: "GetListAccount | get user from token",
		})
	}

	if utils.IsAdmin(userMe.AccountRole) != true {
		return ctx.Status(http.StatusForbidden).JSON(global_types.IResponseAPI{
			Message:      "No Permission",
			ErrorSection: "GetListAccount | validate role user",
		})
	}

	var account []global_types.IAccount

	if err := databases.DB.
		Scopes(services.DebugMode).
		Preload("Logs_Authentication").
		Preload("Creator", services.SelectAccount).
		Preload("Updater", services.SelectAccount).
		Omit("account_password, account_identify_number").
		Find(&account).Error; err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
			Message:      err.Error(),
			ErrorSection: "GetListAccount | find account",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(global_types.IResponseAPI{
		Message: "Get List Account Successfully",
		Data:    account,
	})
}

func GetInformationAccount(ctx *fiber.Ctx) error {
	userMe, err := services.GetContextUser(ctx)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
			Message:      err.Error(),
			ErrorSection: "GetInformationAccount | get user from token",
		})
	}

	body, err := services.GetContextPayload[models.Account](ctx)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
			Message:      "Invalid Payload",
			ErrorSection: "GetInformationAccount | get payload",
		})
	}

	if body.AccountID == 0 {
		return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
			Message:      "Incorrect Parameter",
			ErrorSection: "GetInformationAccount | validate payload",
		})
	}

	if body.AccountID != userMe.AccountID && utils.IsAdmin(userMe.AccountRole) != true {
		return ctx.Status(http.StatusForbidden).JSON(global_types.IResponseAPI{
			Message:      "No Permission",
			ErrorSection: "GetInformationAccount | validate role user",
		})
	}

	var account global_types.IAccount

	var query *gorm.DB = databases.DB.
		Scopes(services.DebugMode).
		Where("account_id = ?", body.AccountID).
		Preload("Logs_Authentication").
		Preload("Creator", services.SelectAccount).
		Preload("Updater", services.SelectAccount).
		Omit("account_password, account_identify_number").
		Find(&account)

	if query.Error != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
			Message:      query.Error.Error(),
			ErrorSection: "GetInformationAccount | find account",
		})
	}

	if query.RowsAffected == 0 {
		return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
			Message:      "Account is not exist",
			ErrorSection: "GetInformationAccount | check account exist",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(global_types.IResponseAPI{
		Message: "Get Information Account Successfully",
		Data:    account,
	})
}

func CreateAccount(ctx *fiber.Ctx) error {
	application, err := services.FindApplicationModel()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
			Message:      err.Error(),
			ErrorSection: "CreateAccount | find application",
		})
	}

	if err := services.IsApplicationUpdatingService(&application); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
			Message:      err.Error(),
			ErrorSection: "CreateAccount | check application update schedule",
		})
	}

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

	if strings.TrimSpace(body.IdentifyNumber) == "" ||
		strings.TrimSpace(body.Username) == "" ||
		strings.TrimSpace(body.Password) == "" ||
		strings.TrimSpace(body.Email) == "" ||
		strings.TrimSpace(body.PhoneNumber) == "" ||
		strings.TrimSpace(body.FirstName) == "" ||
		strings.TrimSpace(body.LastName) == "" {
		return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
			Message:      "Incorrect Parameter",
			ErrorSection: "CreateAccount | validate payload",
		})
	}

	_, err = strconv.Atoi(strings.TrimSpace(body.PhoneNumber))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
			Message:      "Phone number must be a numeric",
			ErrorSection: "CreateAccount | validate format phone number",
		})
	}

	if utils.IsCorrectFormatEmail(body.Email) != true {
		return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
			Message:      "Email is invalid type",
			ErrorSection: "CreateAccount | validate format email",
		})
	}

	if utils.IsCorrectFormatPassword(body.Password) != true {
		return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
			Message:      "Password must contain at least one uppercase and one lowercase and one numeric character",
			ErrorSection: "CreateAccount | validate format password",
		})
	}

	password, err := utils.HashPassword(body.Password)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
			Message:      err.Error(),
			ErrorSection: "CreateAccount | hashing password",
		})
	}

	var duplicate_account models.Account

	var query *gorm.DB = databases.DB.
		Scopes(services.DebugMode).
		Where("account_username = ? OR account_identify_number = ?", strings.TrimSpace(body.Username), strings.TrimSpace(body.IdentifyNumber)).
		Select("account_id").
		Find(&duplicate_account)

	if query.Error != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
			Message:      query.Error.Error(),
			ErrorSection: "CreateAccount | find duplicate account",
		})
	}

	if query.RowsAffected != 0 {
		return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
			Message:      "Duplicate Account",
			ErrorSection: "CreateAccount | check duplicate account",
		})
	}

	var timestamp time.Time = time.Now()

	var account models.Account = models.Account{
		AccountNumber:         "",
		AccountRole:           constants.AccountRoleGeneralUser,
		IdentifyNumber:        strings.TrimSpace(body.IdentifyNumber),
		Username:              strings.TrimSpace(body.Username),
		Password:              password,
		Email:                 strings.TrimSpace(body.Email),
		IsActive:              true,
		IsVerifiedEmail:       false,
		IsVerifiedPhoneNumber: false,
		PhoneNumber:           strings.TrimSpace(body.PhoneNumber),
		FirstName:             strings.TrimSpace(body.FirstName),
		LastName:              strings.TrimSpace(body.LastName),
		CreatedAt:             &timestamp,
	}

	//! *******************************************
	//! ********* TRANSACTION DATABASE ************
	//! *******************************************

	var transaction *gorm.DB = databases.DB.Begin()

	if err := transaction.Scopes(services.DebugMode).Create(&account).Error; err != nil {
		transaction.Rollback()
		return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
			Message:      err.Error(),
			ErrorSection: "CreateAccount | create account",
		})
	}

	if err := transaction.
		Scopes(services.DebugMode).
		Model(&models.Account{}).
		Where("account_id = ?", account.AccountID).
		Update("account_number", fmt.Sprintf("AC%06d", account.AccountID)).
		Update("created_by", account.AccountID).
		Error; err != nil {
		transaction.Rollback()
		return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
			Message:      err.Error(),
			ErrorSection: "CreateAccount | update account number",
		})
	}

	transaction.Commit()

	return ctx.Status(fiber.StatusOK).JSON(global_types.IResponseAPI{
		Message: "Create Account Successfully",
	})
}
