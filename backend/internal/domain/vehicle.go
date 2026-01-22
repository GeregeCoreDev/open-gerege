// Package domain provides implementation for domain
//
// File: vehicle.go
// Description: implementation for domain
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package domain

type Vehicle struct {
	Id      int    `json:"id" gorm:"primaryKey"`
	PlateNo string `json:"plate_no" gorm:"unique;not null;size:7"`
}
