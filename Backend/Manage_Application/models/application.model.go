package models

import (
	"time"
)

type Application struct {
	ApplicationID                      int        `json:"application_id,omitempty" gorm:"primaryKey; autoIncrement; column:application_id; comment:Primary Key"`
	ApplicationName                    string     `json:"application_name,omitempty" gorm:"type:varchar(255); column:application_name;"`
	ApplicationScheduleUpdateStartTime *time.Time `json:"application_schedule_update_start_time,omitempty" gorm:"type:timestamp; not null; column:application_schedule_update_start_time;"`
	ApplicationScheduleUpdateEndTime   *time.Time `json:"application_schedule_update_end_time,omitempty" gorm:"type:timestamp; not null; column:application_schedule_update_end_time;"`
	CreatedAt                          *time.Time `json:"created_at,omitempty" gorm:"type:timestamp; not null; column:created_at;"`
	CreatedBy                          int        `json:"created_by,omitempty" gorm:"type:bigint; not null; column:created_by;"`
	UpdatedAt                          *time.Time `json:"updated_at,omitempty" gorm:"type:timestamp; default:null; column:updated_at;"`
	UpdatedBy                          int        `json:"updated_by,omitempty" gorm:"type:bigint; default:null; column:updated_by;"`
}

func (Application *Application) TableName() string {
	return "application"
}
