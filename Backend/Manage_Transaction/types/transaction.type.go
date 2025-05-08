package global_types

import "chatify/models"

type ITransaction struct {
	models.Transaction
	ReferenceTransaction models.Transaction       `json:"reference_transaction" gorm:"foreignKey:TransactionID; references:ReferenceTransactionID"`
	TransactionFile      []models.TransactionFile `json:"transaction_file" gorm:"foreignKey:TransactionID"`
	TransactionReference []ITransaction           `json:"transaction_reference" gorm:"foreignKey:ReferenceTransactionID; references:TransactionID"`
	Creator              models.Account           `json:"creator" gorm:"foreignKey:AccountID; references:CreatedBy"`
	Updater              models.Account           `json:"updater" gorm:"foreignKey:AccountID; references:UpdatedBy"`
}
