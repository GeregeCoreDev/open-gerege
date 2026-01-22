// Package dto provides implementation for dto
//
// File: notification_dto.go
// Description: implementation for dto
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package dto

type NotificationReadDto struct {
	GroupId int `json:"group_id" validate:"required,gt=0"`
}

type NotificationSendDto struct {
	Tenant        string `json:"tenant" validate:"required"`
	UserID        int    `json:"user_id"` // 0 бол broadcast_all
	Title         string `json:"title"`
	Content       string `json:"content"`
	IdempotentKey string `json:"idempotency_key"`
}
