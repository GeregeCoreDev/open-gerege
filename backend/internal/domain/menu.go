package domain

type Menu struct {
	ID           int64       `json:"id" gorm:"primaryKey"`
	Code         string      `json:"code" gorm:"type:varchar(255)"`
	Key          string      `json:"key" gorm:"type:varchar(255);unique"`
	Name         string      `json:"name" gorm:"type:varchar(255)"`
	Description  string      `json:"description" gorm:"type:varchar(255)"`
	Icon         string      `json:"icon" gorm:"type:varchar(255)"`
	Path         string      `json:"path" gorm:"type:varchar(255)"`
	Sequence     int64       `json:"sequence"`
	ParentID     *int64      `json:"parent_id"`
	Parent       *Menu       `json:"parent,omitempty" gorm:"foreignKey:ParentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Children     []Menu      `json:"children,omitempty" gorm:"foreignKey:ParentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PermissionID *int64      `json:"permission_id"`
	Permission   *Permission `json:"permission,omitempty" gorm:"foreignKey:PermissionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	IsActive     *bool       `json:"is_active"`
	ExtraFields
}
