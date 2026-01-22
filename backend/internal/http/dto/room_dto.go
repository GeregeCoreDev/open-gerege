// Package dto provides implementation for dto
//
// File: room_dto.go
// Description: implementation for dto
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package dto

import "time"

type RoomResponse struct {
	ID          int        `json:"id"`
	UUID        string     `json:"uuid"`
	Name        string     `json:"name"`
	Code        string     `json:"code"`
	MeetingDate *time.Time `json:"meeting_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   time.Time  `json:"expires_at"`
	UserCount   int        `json:"user_count"`
	CreatedByID int        `json:"created_by_user_id"`
	IsCreator   bool       `json:"is_creator"`
}

type CreateRoomRequest struct {
	Name        string     `json:"name" validate:"required,min=3,max=100"`
	MeetingDate *time.Time `json:"meeting_date" validate:"omitempty"`
	UserIDs     []int      `json:"user_ids"`
}

// AddUsersRequest represents the request body for adding users to a room
type AddUsersRequest struct {
	UserIDs []int `json:"user_ids" validate:"required,min=1,dive,gt=0"`
}

// JoinRoomRequest represents the request body for joining a room by code
type JoinRoomRequest struct {
	Code string `json:"code" validate:"required,len=4,numeric"`
}
