// Package domain provides implementation for domain
//
// File: system.go
// Description: implementation for domain
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package domain

type System struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	Code        string `json:"code" gorm:"type:varchar(255);unique"`
	Key         string `json:"key" gorm:"type:varchar(255)"`
	Name        string `json:"name" gorm:"type:varchar(255)"`
	Description string `json:"description" gorm:"type:varchar(255)"`
	IsActive    *bool  `json:"is_active"`
	Icon        string `json:"icon" gorm:"type:varchar(255)"`
	Sequence    int    `json:"sequence"`
	ExtraFields
}
