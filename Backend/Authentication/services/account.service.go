package services

import (
	"chatify/databases"
	"chatify/models"
	"errors"

	"gorm.io/gorm"
)

func FindAccountModel(account_id *int) (models.Account, error) {
	var account models.Account

	query := databases.DB.
		Scopes(DebugMode).
		Where("is_active = true").
		Where("account_id = ?", account_id).
		Scopes(SelectAccount).
		Find(&account)

	if query.Error != nil {
		return account, query.Error
	}

	if query.RowsAffected == 0 {
		return account, errors.New("Account is not exist")
	}

	return account, nil
}

func SelectAccount(db *gorm.DB) *gorm.DB {
	return db.Unscoped().Omit("account_password, account_identify_number").Select("account_id, account_role, account_first_name, account_last_name")
}
