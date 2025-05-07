package global_types

type IFile struct {
	FileName string `json:"file_name,omitempty" gorm:"type:varchar(50); column:file_name;"`
	FilePath string `json:"file_path,omitempty" gorm:"type:varchar(50); column:file_path;"`
	FileType string `json:"file_type,omitempty" gorm:"type:varchar(50); column:file_type;"`
}
