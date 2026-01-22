// Package domain provides implementation for domain
//
// File: permission.go
// Description: implementation for domain
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package domain

type Permission struct {
	ID          int     `json:"id" gorm:"primaryKey"`
	Code        string  `json:"code" gorm:"unique;not null;type:varchar(255)"`
	Name        string  `json:"name" gorm:"type:varchar(255)"`
	Description string  `json:"description" gorm:"type:varchar(255)"`
	SystemID    int     `json:"system_id"`
	System      *System `json:"system,omitempty" gorm:"foreignKey:SystemID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ModuleID    int     `json:"module_id"`
	Module      *Module `json:"module,omitempty" gorm:"foreignKey:ModuleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ActionID    *int64  `json:"action_id"`
	Action      *Action `json:"action,omitempty" gorm:"foreignKey:ActionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	IsActive    *bool   `json:"is_active" gorm:"not null;default:true"`
	ExtraFields
}
