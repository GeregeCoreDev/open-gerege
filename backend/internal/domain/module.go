// Package domain provides implementation for domain
//
// File: module.go
// Description: implementation for domain
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package domain

type Module struct {
	ID          int     `json:"id" gorm:"primaryKey"`
	Code        string  `json:"code" gorm:"type:varchar(255);unique"`
	Name        string  `json:"name" gorm:"type:varchar(255)"`
	Description string  `json:"description" gorm:"type:varchar(255)"`
	IsActive    *bool   `json:"is_active"`
	SystemID    int     `json:"system_id"`
	System      *System `json:"system,omitempty" gorm:"foreignKey:SystemID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ExtraFields
}
