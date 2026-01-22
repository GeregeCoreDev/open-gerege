// Package domain provides implementation for domain
//
// File: news.go
// Description: implementation for domain
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package domain

type News struct {
	Id       int    `json:"id" gorm:"primaryKey"`
	Title    string `json:"title" gorm:"type:varchar(255)"`
	Text     string `json:"text" gorm:"type:text"`
	ImageUrl string `json:"image_url" gorm:"type:varchar(255)"`
	ExtraFields
}
