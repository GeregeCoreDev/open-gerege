// Package domain provides implementation for domain
//
// File: role.go
// Description: implementation for domain
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package domain

type Role struct {
	ID           int     `json:"id" gorm:"primaryKey"`
	SystemID     int     `json:"system_id"`
	System       *System `json:"system,omitempty" gorm:"foreignKey:SystemID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Code         string  `json:"code" gorm:"unique;not null;type:varchar(255)"`
	Name         string  `json:"name" gorm:"not null;type:varchar(255)"`
	Description  string  `json:"description" gorm:"type:varchar(255)"`
	IsActive     *bool   `json:"is_active"`
	IsSystemRole *bool   `json:"is_system_role" gorm:"default:false"`
	ExtraFields
}

type RolePermission struct {
	RoleID       int         `json:"role_id"`
	PermissionID int         `json:"permission_id"`
	Permission   *Permission `json:"permission,omitempty" gorm:"foreignKey:PermissionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ExtraFields
}
