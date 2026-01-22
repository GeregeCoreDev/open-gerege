// Package dto provides implementation for dto
//
// File: terminal_dto.go
// Description: implementation for dto
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package dto

type TerminalCreateDto struct {
	Serial string `json:"serial" validate:"required,max=80"`
	Name   string `json:"name"   validate:"required,max=255"`
	OrgId  *int   `json:"org_id,omitempty" validate:"omitempty,gt=0"` // хоосон бол token-оос авна
}

type TerminalUpdateDto TerminalCreateDto
