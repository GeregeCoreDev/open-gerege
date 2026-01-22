// Package domain provides implementation for domain
//
// File: notification.go
// Description: implementation for domain
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package domain

type NotificationGroup struct {
	Id              int    `gorm:"primaryKey" json:"id"`
	UserId          int    `json:"user_id" gorm:"index"`
	Title           string `json:"title" gorm:"size:255"`
	Content         string `json:"content" gorm:"type:text"`
	Type            string `json:"type" gorm:"size:20"`
	Tenant          string `json:"tenant" gorm:"size:50"`
	CreatedUsername string `json:"created_username" gorm:"size:100"`
	ExtraFields
}

type Notification struct {
	Id              int    `gorm:"primaryKey" json:"id"`
	UserId          int    `json:"user_id" gorm:"index"`
	Title           string `json:"title" gorm:"size:255"`
	Content         string `json:"content" gorm:"type:text"`
	IsRead          bool   `json:"is_read" gorm:"default:false"`
	Type            string `json:"type" gorm:"size:20"`
	Tenant          string `json:"tenant" gorm:"size:50"`
	GroupId         int    `json:"group_id" gorm:"index"`
	CreatedUsername string `json:"created_username" gorm:"size:100"`
	ExtraFields
}
