// Package domain provides implementation for domain
//
// File: terminal.go
// Description: implementation for domain
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package domain

type Terminal struct {
	Id           int           `json:"id" gorm:"primaryKey"`
	Serial       string        `json:"serial" gorm:"unique;not null;type:varchar(80)" validate:"required"`
	Name         string        `json:"name" gorm:"type:varchar(255)" validate:"required"`
	OrgId        *int          `json:"org_id"`
	Organization *Organization `json:"organization,omitempty" gorm:"foreignKey:OrgId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ExtraFields
}
