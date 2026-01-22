// Package dto provides implementation for dto
//
// File: user_role_dto.go
// Description: implementation for dto
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package dto

import "git.gerege.mn/backend-packages/common"

type UserRoleUsersQuery struct {
	RoleID int `query:"role_id" validate:"required"`
	common.PaginationQuery
}

type UserRoleRolesQuery struct {
	UserID int `query:"user_id" validate:"required"`
	common.PaginationQuery
}

type UserRoleAssignByRole struct {
	RoleID  int   `json:"role_id"  validate:"required"`
	UserIDs []int `json:"user_ids" validate:"required,min=1,dive,gt=0"`
}

type UserRoleAssignByUser struct {
	UserID  int   `json:"user_id"  validate:"required"`
	RoleIDs []int `json:"role_ids" validate:"required,min=1,dive,gt=0"`
}

type UserRoleRemoveDto struct {
	UserID int `json:"user_id" validate:"required"`
	RoleID int `json:"role_id" validate:"required"`
}
