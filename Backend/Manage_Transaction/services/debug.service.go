package services

import (
	"chatify/configs"

	"gorm.io/gorm"
)

func DebugMode(db *gorm.DB) *gorm.DB {
	if configs.ENV.IsProductionMode != true {
		return db.Debug()
	}
	return db
}
