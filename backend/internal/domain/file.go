// Package domain provides implementation for domain
//
// File: file.go
// Description: implementation for domain
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package domain

type PublicFile struct {
	Id          int    `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"type:varchar(255)"`
	Extension   string `json:"extension" gorm:"type:varchar(10)"`
	Description string `json:"description" gorm:"type:varchar(255)"`
	FileUrl     string `json:"file_url" gorm:"type:varchar(255)"`
	ExtraFields
}
