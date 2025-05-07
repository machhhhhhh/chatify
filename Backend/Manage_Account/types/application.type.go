package global_types

import "chatify/models"

type IApplication struct {
	models.Application
	ApplicationIsAllowToPerformTaskWorkload bool           `json:"application_is_allow_to_perform_task_workload"`
	Creator                                 models.Account `json:"creator" gorm:"foreignKey:AccountID; references:CreatedBy"`
	Updater                                 models.Account `json:"updater" gorm:"foreignKey:AccountID; references:UpdatedBy"`
}
