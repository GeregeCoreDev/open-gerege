package domain

type Action struct {
	ID          int64  `json:"id" gorm:"primaryKey"`
	Code        string `json:"code" gorm:"unique;not null;type:varchar(255)"`
	Name        string `json:"name" gorm:"type:varchar(255)"`
	Description string `json:"description" gorm:"type:varchar(255)"`
	IsActive    *bool  `json:"is_active" gorm:"not null;default:true"`
	ExtraFields
}
