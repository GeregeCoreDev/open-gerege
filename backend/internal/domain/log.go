// Package domain provides implementation for domain
//
// File: log.go
// Description: implementation for domain
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package domain

import (
	"time"

	"gorm.io/datatypes"
)

type APILog struct {
	Id          int64          `json:"id" gorm:"primaryKey"`
	OrgId       *int64         `json:"org_id" gorm:"index"`
	UserId      *int64         `json:"user_id" gorm:"index"`
	Username    string         `json:"username" gorm:"type:varchar(50)"`
	Path        string         `json:"path" gorm:"type:varchar(255)"`
	Method      string         `json:"method" gorm:"type:varchar(10)"`
	Params      datatypes.JSON `json:"-" gorm:"type:jsonb"`
	Queries     datatypes.JSON `json:"-" gorm:"type:jsonb"`
	Body        datatypes.JSON `json:"-" gorm:"type:jsonb"`
	StatusCode  int            `json:"status_code"`
	Response    datatypes.JSON `json:"-" gorm:"type:jsonb"`
	LatencyMs   int64          `gorm:"column:latency_ms"`
	ReqSize     int64          `gorm:"column:req_size"`
	ResSize     int64          `gorm:"column:res_size"`
	IP          string         `gorm:"size:45;column:ip"`
	CreatedDate time.Time      `json:"created_date" gorm:"column:created_date"`
}

// TableName specifies the table name for APILog
func (APILog) TableName() string {
	return "template_backend.logs"
}

type AuditLog struct {
	ID          uint    `gorm:"primaryKey"`
	UserID      *int    `gorm:"index"`
	ActionID    int64   `gorm:"index"`
	Action      *Action `gorm:"foreignKey:ActionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ModuleID    int64   `gorm:"index"`
	Module      *Module `gorm:"foreignKey:ModuleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Metadata    datatypes.JSON
	CreatedDate time.Time `json:"created_date"`
}
