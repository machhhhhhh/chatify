package global_types

import "chatify/models"

type IAccount struct {
	models.Account
	Logs_Authentication []models.Logs_Authentication `json:"logs_authentication" gorm:"foreignKey:AccountID"`
	Creator             models.Account               `json:"creator" gorm:"foreignKey:AccountID; references:CreatedBy"`
	Updater             models.Account               `json:"updater" gorm:"foreignKey:AccountID; references:UpdatedBy"`
}
