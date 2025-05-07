package models

import (
	"chatify/constants"
	"time"
)

type Logs_Authentication struct {
	LogsAuthenticationID    int                                  `json:"logs_authentication_id,omitempty" gorm:"primaryKey; autoIncrement; column:logs_authentication_id; comment:Primary Key"`
	AccountID               int                                  `json:"account_id,omitempty" gorm:"type:bigint; column:account_id;"`
	LogsActivity            constants.LogsActivityAuthentication `json:"logs_activity,omitempty" gorm:"type:text; column:logs_activity; default:'LOGIN'; comment:'LOGIN, LOGOUT'"`
	ReceiveNotificationAt   *time.Time                           `json:"receive_notification_at,omitempty" gorm:"type:timestamp; column:receive_notification_at;"`
	LogsNotificationCounter int                                  `json:"logs_notification_counter,omitempty" gorm:"type:bigint; column:logs_notification_counter;"`
	CreatedAt               *time.Time                           `json:"created_at,omitempty" gorm:"type:timestamp; not null; column:created_at;"`
	CreatedBy               int                                  `json:"created_by,omitempty" gorm:"type:bigint; not null; column:created_by;"`
	ExpiredAt               *time.Time                           `json:"expired_at,omitempty" gorm:"type:timestamp; column:expired_at;"`
}

func (Logs_Authentication *Logs_Authentication) TableName() string {
	return "logs_authentication"
}
