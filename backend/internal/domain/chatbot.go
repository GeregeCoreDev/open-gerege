// Package domain provides implementation for domain
//
// File: chatbot.go
// Description: implementation for domain
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package domain

type ChatItem struct {
	ID     int    `json:"id"`
	Key    string `json:"key"`
	Answer string `json:"answer"`
	ExtraFields
}
