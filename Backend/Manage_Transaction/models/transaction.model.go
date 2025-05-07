package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	TransactionID          int            `json:"transaction_id,omitempty" gorm:"primaryKey; autoIncrement; column:transaction_id; comment:Primary Key"`
	ReferenceTransactionID int            `json:"reference_transaction_id,omitempty" gorm:"type:bigint; not null; column:reference_transaction_id;"`
	TransactionNumber      string         `json:"transaction_number,omitempty" gorm:"type:varchar(50); column:transaction_number;"`
	TransactionDescription string         `json:"transction_description,omitempty" gorm:"type:varchar(255); column:transction_description;"`
	IsActive               bool           `json:"is_active" gorm:"type:boolean; column:is_active;"`
	CreatedAt              *time.Time     `json:"created_at,omitempty" gorm:"type:timestamp; not null; column:created_at;"`
	CreatedBy              int            `json:"created_by,omitempty" gorm:"type:bigint; not null; column:created_by;"`
	UpdatedAt              *time.Time     `json:"updated_at,omitempty" gorm:"type:timestamp; default:null; column:updated_at;"`
	UpdatedBy              int            `json:"updated_by,omitempty" gorm:"type:bigint; default:null; column:updated_by;"`
	DeletedAt              gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"type:timestamp; column:deleted_at;" swaggerignore:"true"`
	DeletedBy              int            `json:"deleted_by,omitempty" gorm:"type:bigint; default:null; column:deleted_by;" swaggerignore:"true"`
}

func (Transaction *Transaction) TableName() string {
	return "transaction"
}
