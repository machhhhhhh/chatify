package models

import (
	"time"

	"gorm.io/gorm"
)

type TransactionFile struct {
	TransactionFileID   int            `json:"transaction_file_id,omitempty" gorm:"primaryKey; autoIncrement; column:transaction_file_id; comment:Primary Key"`
	TransactionID       int            `json:"transaction_id,omitempty" gorm:"type:bigint; not null; column:transaction_id;"`
	TransactionFileName string         `json:"transaction_file_name,omitempty" gorm:"type:varchar(50); column:transaction_file_name;"`
	TransactionFilePath string         `json:"transaction_file_path,omitempty" gorm:"type:varchar(50); column:transaction_file_path;"`
	TransactionFileType string         `json:"transaction_file_type,omitempty" gorm:"type:varchar(50); column:transaction_file_type;"`
	IsActive            bool           `json:"is_active" gorm:"type:boolean; column:is_active;"`
	CreatedAt           *time.Time     `json:"created_at,omitempty" gorm:"type:timestamp; not null; column:created_at;"`
	CreatedBy           int            `json:"created_by,omitempty" gorm:"type:bigint; not null; column:created_by;"`
	UpdatedAt           *time.Time     `json:"updated_at,omitempty" gorm:"type:timestamp; default:null; column:updated_at;"`
	UpdatedBy           int            `json:"updated_by,omitempty" gorm:"type:bigint; default:null; column:updated_by;"`
	DeletedAt           gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"type:timestamp; column:deleted_at;" swaggerignore:"true"`
	DeletedBy           int            `json:"deleted_by,omitempty" gorm:"type:bigint; default:null; column:deleted_by;" swaggerignore:"true"`
}

func (TransactionFile *TransactionFile) TableName() string {
	return "transaction_file"
}
