package controllers

import (
	"chatify/configs"
	"chatify/databases"
	"chatify/models"
	"chatify/services"
	global_types "chatify/types"
	"chatify/utils"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func MigrateTransaction(ctx *fiber.Ctx) error {
	// userMe, err := services.GetContextUser(ctx)
	// if err != nil {
	// 	return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
	// 		Message:      err.Error(),
	// 		ErrorSection: "MigrateTransaction | get user from token",
	// 	})
	// }

	// if utils.IsAdmin(userMe.AccountRole) != true {
	// 	return ctx.Status(http.StatusForbidden).JSON(global_types.IResponseAPI{
	// 		Message:      "No Permission",
	// 		ErrorSection: "MigrateTransaction | validate role user",
	// 	})
	// }

	// application, err := services.FindApplicationModel()
	// if err != nil {
	// 	return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
	// 		Message:      err.Error(),
	// 		ErrorSection: "MigrateTransaction | find application",
	// 	})
	// }

	// if err := services.IsApplicationUpdatingService(&application); err != nil {
	// 	return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
	// 		Message:      err.Error(),
	// 		ErrorSection: "MigrateTransaction | check application update schedule",
	// 	})
	// }

	//! *******************************************
	//! ********* TRANSACTION DATABASE ************
	//! *******************************************

	var transaction *gorm.DB = databases.DB.Begin()

	if err := transaction.Scopes(services.DebugMode).AutoMigrate(
		&models.Transaction{},
		&models.TransactionFile{},
	); err != nil {
		transaction.Rollback()
		return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
			Message:      err.Error(),
			ErrorSection: "MigrateTransaction | migrate data",
		})
	}

	transaction.Commit()

	return ctx.Status(fiber.StatusOK).JSON(global_types.IResponseAPI{
		Message: "Migrate Transaction Successfully",
	})
}

func GetListTransaction(ctx *fiber.Ctx) error {
	type IRequestBody struct {
		TransactionID int  `json:"number"`
		IsShowComment bool `json:"is_show_comment"`
	}

	body, err := services.GetContextPayload[IRequestBody](ctx)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
			Message:      "Invalid Payload",
			ErrorSection: "GetListTransaction | get payload",
		})
	}

	if body.TransactionID < 0 {
		return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
			Message:      "Incorrect Parameter",
			ErrorSection: "CreateTransaction | validate payload",
		})
	}

	var transaction []global_types.ITransaction

	if err := databases.DB.
		Scopes(services.DebugMode).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if body.TransactionID != 0 {
				return db.Where("reference_transaction_id = ?", body.TransactionID)
			}
			return db
		}).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if body.IsShowComment != true {
				return db.Where("CASE WHEN reference_transaction_id IS NOT NULL AND reference_transaction_id != 0 THEN FALSE ELSE TRUE END")
			}
			return db
		}).
		Preload("TransactionFile", func(db *gorm.DB) *gorm.DB {
			return db.Select("*, CONCAT(?::text,'/',transaction_file_path) as transaction_file_path", configs.ENV.FileSetting.PathRenderFile)
		}).
		Preload("TransactionReference").
		Preload("TransactionReference.Creator", services.SelectAccount).
		Preload("TransactionReference.TransactionFile", func(db *gorm.DB) *gorm.DB {
			return db.Select("*, CONCAT(?::text,'/',transaction_file_path) as transaction_file_path", configs.ENV.FileSetting.PathRenderFile)
		}).
		Preload("Creator", services.SelectAccount).
		Preload("Updater", services.SelectAccount).
		Order("transaction_id DESC").
		Omit("deleted_at").
		Find(&transaction).Error; err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
			Message:      err.Error(),
			ErrorSection: "GetListTransaction | find transaction",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(global_types.IResponseAPI{
		Message: "Get List Tranasction Successfully",
		Data:    transaction,
	})
}

func GetInformationTransaction(ctx *fiber.Ctx) error {

	body, err := services.GetContextPayload[models.Transaction](ctx)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
			Message:      "Invalid Payload",
			ErrorSection: "GetInformationTransaction | get payload",
		})
	}

	if body.TransactionID <= 0 {
		return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
			Message:      "Incorrect Parameter",
			ErrorSection: "GetInformationTransaction | validate payload",
		})
	}

	var transaction global_types.ITransaction

	var query *gorm.DB = databases.DB.
		Scopes(services.DebugMode).
		Where("transaction_id = ?", body.TransactionID).
		Preload("TransactionFile", func(db *gorm.DB) *gorm.DB {
			return db.Select("*, CONCAT(?::text,'/',transaction_file_path) as transaction_file_path", configs.ENV.FileSetting.PathRenderFile)
		}).
		Preload("ReferenceTransaction").
		Preload("TransactionReference").
		Preload("TransactionReference.Creator", services.SelectAccount).
		Preload("TransactionReference.TransactionFile", func(db *gorm.DB) *gorm.DB {
			return db.Select("*, CONCAT(?::text,'/',transaction_file_path) as transaction_file_path", configs.ENV.FileSetting.PathRenderFile)
		}).
		Preload("Creator", services.SelectAccount).
		Preload("Updater", services.SelectAccount).
		Find(&transaction)

	if query.Error != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
			Message:      query.Error.Error(),
			ErrorSection: "GetInformationTransaction | find transaction",
		})
	}

	if query.RowsAffected == 0 {
		return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
			Message:      "Transaction is not exist",
			ErrorSection: "GetInformationTransaction | check transaction exist",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(global_types.IResponseAPI{
		Message: "Get Information Transaction Successfully",
		Data:    transaction,
	})
}

func CreateTransaction(ctx *fiber.Ctx) error {
	userMe, err := services.GetContextUser(ctx)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
			Message:      err.Error(),
			ErrorSection: "CreateTransaction | get user from token",
		})
	}

	application, err := services.FindApplicationModel()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
			Message:      err.Error(),
			ErrorSection: "CreateTransaction | find application",
		})
	}

	if err := services.IsApplicationUpdatingService(&application); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
			Message:      err.Error(),
			ErrorSection: "CreateTransaction | check application update schedule",
		})
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
			Message:      "Error parsing multipart form",
			ErrorSection: "Authorization | parsing multipart/form-data",
		})
	}

	var request_file []*multipart.FileHeader = form.File["files"]

	body, err := services.GetContextPayload[global_types.ITransaction](ctx)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
			Message:      "Invalid Payload",
			ErrorSection: "CreateTransaction | get payload",
		})
	}

	if body.TransactionID < 0 || len(strings.TrimSpace(body.TransactionDescription)) == 0 {
		return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
			Message:      "Incorrect Parameter",
			ErrorSection: "CreateTransaction | validate payload",
		})
	}

	if body.TransactionID != 0 {
		var transaction models.Transaction

		var query *gorm.DB = databases.DB.
			Scopes(services.DebugMode).
			Where("transaction_id = ?", body.TransactionID).
			Select("transaction_id").
			Find(&transaction)

		if query.Error != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
				Message:      query.Error.Error(),
				ErrorSection: "CreateTransaction | find transaction",
			})
		}

		if query.RowsAffected == 0 {
			return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
				Message:      "Transaction is not exist",
				ErrorSection: "CreateTransaction | check transaction exist",
			})
		}
	}

	//! *******************************************
	//! ********* TRANSACTION DATABASE ************
	//! *******************************************

	var transaction *gorm.DB = databases.DB.Begin()

	var timestamp time.Time = time.Now()

	var transaction_model models.Transaction = models.Transaction{
		TransactionNumber:      "",
		ReferenceTransactionID: body.TransactionID,
		TransactionDescription: strings.TrimSpace(body.TransactionDescription),
		IsActive:               true,
		CreatedAt:              &timestamp,
		CreatedBy:              userMe.AccountID,
	}

	if err := transaction.Scopes(services.DebugMode).Create(&transaction_model).Error; err != nil {
		transaction.Rollback()
		return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
			Message:      err.Error(),
			ErrorSection: "CreateTransaction | create transaction",
		})
	}

	if err := transaction.
		Scopes(services.DebugMode).
		Model(&models.Transaction{}).
		Where("transaction_id = ?", transaction_model.TransactionID).
		Update("transaction_number", fmt.Sprintf("TRANXID%06d", transaction_model.TransactionID)).
		Error; err != nil {
		transaction.Rollback()
		return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
			Message:      err.Error(),
			ErrorSection: "CreateTransaction | update transaction number",
		})
	}

	var transaction_file []models.TransactionFile

	if len(request_file) != 0 {
		if err := services.CreateDirectory(configs.ENV.FileSetting.RootDirectory); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
				Message:      "Cannot create directory",
				ErrorSection: "CreateTransaction | create directory",
			})
		}

		for i := range request_file {
			var extension string = filepath.Ext(request_file[i].Filename)
			var file_name string = uuid.New().String() + extension
			var file_path string = utils.GetFileDirectory(file_name)

			transaction_file = append(transaction_file, models.TransactionFile{
				TransactionID:       transaction_model.TransactionID,
				TransactionFileName: request_file[i].Filename,
				TransactionFilePath: file_name,
				TransactionFileType: request_file[i].Header["Content-Type"][0],
				IsActive:            true,
				CreatedAt:           &timestamp,
				CreatedBy:           userMe.AccountID,
			})

			err := ctx.SaveFile(request_file[i], file_path)
			if err != nil {
				transaction.Rollback()
				return ctx.Status(http.StatusBadRequest).JSON(global_types.IResponseAPI{
					Message:      "Cannot save file",
					ErrorSection: "CreateTransaction | save file",
				})
			}
		}

		if err := transaction.Scopes(services.DebugMode).Create(&transaction_file).Error; err != nil {
			transaction.Rollback()
			return ctx.Status(http.StatusInternalServerError).JSON(global_types.IResponseAPI{
				Message:      err.Error(),
				ErrorSection: "CreateTransaction | create transaction file",
			})
		}
	}

	transaction.Commit()

	//! *******************************************
	//! ********** EXTERNAL SERVICE ***************
	//! *******************************************

	// TODO: send socket-client to socket-server
	// TODO: for triggering the client-side terms of event socket-room

	return ctx.Status(fiber.StatusOK).JSON(global_types.IResponseAPI{
		Message: "Create Transaction Successfully",
		Data:    transaction_model.TransactionID,
	})
}
