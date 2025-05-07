package services

import (
	"chatify/databases"
	global_types "chatify/types"
	"errors"

	"gorm.io/gorm"
)

func IsApplicationUpdatingService(application *global_types.IApplication) error {
	if application.ApplicationIsAllowToPerformTaskWorkload == true {
		return nil
	}
	return errors.New("Application is now updating. Please try again later")
}

func FindApplicationModel() (global_types.IApplication, error) {
	var application global_types.IApplication

	query := databases.DB.
		Scopes(DebugMode).
		Scopes(SelectApplication).
		Find(&application)

	if query.Error != nil {
		return application, query.Error
	}

	if query.RowsAffected == 0 {
		return application, errors.New("Application configuration is not exist")
	}

	return application, nil
}

func SelectApplication(db *gorm.DB) *gorm.DB {
	return db.Select(`
	*,
	(CASE 
		WHEN current_time BETWEEN application_schedule_update_start_time::time AND application_schedule_update_end_time::time
		THEN false ELSE true
	END) AS "ApplicationIsAllowToPerformTaskWorkload"
`)
}
