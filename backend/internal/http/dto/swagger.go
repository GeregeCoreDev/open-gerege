// Package dto provides data transfer objects
//
// File: swagger.go
// Description: Swagger documentation models
//
// These types are used for Swagger/OpenAPI documentation only.
// They provide proper type definitions for API responses.
package dto

// Response is the standard API response structure
// @Description Standard API response wrapper
type Response struct {
	// Business logic code (OK, BAD_REQUEST, etc.)
	Code string `json:"code" example:"OK"`
	// Human-readable message
	Message string `json:"message" example:"success"`
	// Unique request identifier
	RequestID string `json:"request_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	// Response data (can be any type)
	Data interface{} `json:"data,omitempty"`
	// Error details for validation errors
	Details interface{} `json:"details,omitempty"`
}

// ErrorResponse is the standard error response
// @Description Error response structure
type ErrorResponse struct {
	// Error code
	Code string `json:"code" example:"BAD_REQUEST"`
	// Error message
	Message string `json:"message" example:"Invalid request parameters"`
	// Request ID for debugging
	RequestID string `json:"request_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	// Detailed error information
	Details map[string]string `json:"details,omitempty"`
}

// PaginatedResponse is the paginated list response
// @Description Paginated response structure
type PaginatedResponse struct {
	// Business logic code
	Code string `json:"code" example:"OK"`
	// Message
	Message string `json:"message" example:"success"`
	// Request ID
	RequestID string `json:"request_id"`
	// Response data
	Data interface{} `json:"data"`
	// Pagination metadata
	Meta PaginationMeta `json:"meta"`
}

// PaginationMeta contains pagination information
// @Description Pagination metadata
type PaginationMeta struct {
	// Current page number
	Page int `json:"page" example:"1"`
	// Page size
	Size int `json:"size" example:"20"`
	// Total number of items
	Total int64 `json:"total" example:"100"`
	// Total number of pages
	TotalPages int `json:"total_pages" example:"5"`
}
