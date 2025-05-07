package models

import (
	"chatify/constants"
	"time"

	"gorm.io/gorm"
)

type Account struct {
	AccountID             int                   `json:"account_id,omitempty" gorm:"primaryKey; autoIncrement; column:account_id; comment:Primary Key"`
	AccountNumber         string                `json:"account_number,omitempty" gorm:"type:varchar(50); column:account_number;"`
	AccountRole           constants.AccountRole `json:"account_role,omitempty" gorm:"type:text; column:account_role; default:'GENERAL_USER'; comment:'GENERAL_USER, EMPLOYEE, ADMIN, SUPER_ADMIN'"`
	IdentifyNumber        string                `json:"account_identify_number,omitempty" gorm:"type:varchar(255); column:account_identify_number;"`
	FirstName             string                `json:"account_first_name,omitempty" gorm:"type:varchar(255); column:account_first_name;"`
	LastName              string                `json:"account_last_name,omitempty" gorm:"type:varchar(255); column:account_last_name;"`
	Username              string                `json:"account_username,omitempty" gorm:"type:varchar(255); column:account_username;"`
	Password              string                `json:"account_password,omitempty" gorm:"type:varchar(255); column:account_password;"`
	Email                 string                `json:"account_email,omitempty" gorm:"type:varchar(255); column:account_email;"`
	PhoneNumber           string                `json:"account_phone_number,omitempty" gorm:"type:varchar(255); column:account_phone_number;"`
	IsVerifiedEmail       bool                  `json:"account_is_verified_email" gorm:"type:boolean; column:account_is_verified_email;"`
	IsVerifiedPhoneNumber bool                  `json:"account_is_verified_phone_number" gorm:"type:boolean; column:account_is_verified_phone_number;"`
	LastLogin             *time.Time            `json:"last_login,omitempty" gorm:"type:timestamp; default:null; column:last_login;"`
	IsActive              bool                  `json:"is_active" gorm:"type:boolean; column:is_active;"`
	CreatedAt             *time.Time            `json:"created_at,omitempty" gorm:"type:timestamp; not null; column:created_at;"`
	CreatedBy             int                   `json:"created_by,omitempty" gorm:"type:bigint; not null; column:created_by;"`
	UpdatedAt             *time.Time            `json:"updated_at,omitempty" gorm:"type:timestamp; default:null; column:updated_at;"`
	UpdatedBy             int                   `json:"updated_by,omitempty" gorm:"type:bigint; default:null; column:updated_by;"`
	DeletedAt             gorm.DeletedAt        `json:"deleted_at,omitempty" gorm:"type:timestamp; column:deleted_at;" swaggerignore:"true"`
	DeletedBy             int                   `json:"deleted_by,omitempty" gorm:"type:bigint; default:null; column:deleted_by;" swaggerignore:"true"`
}

func (Account *Account) TableName() string {
	return "account"
}
