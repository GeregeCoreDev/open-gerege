// Package errors provides custom error types for the application
//
// File: errors.go
// Description: Custom error types for business logic and infrastructure errors
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package errors

import (
	"errors"
	"fmt"
)

// Error codes for API responses
const (
	CodeNotFound       = "NOT_FOUND"
	CodeValidation     = "VALIDATION_ERROR"
	CodeUnauthorized   = "UNAUTHORIZED"
	CodeForbidden      = "FORBIDDEN"
	CodeConflict       = "CONFLICT"
	CodeInternal       = "INTERNAL_ERROR"
	CodeBadRequest     = "BAD_REQUEST"
	CodeExternalAPI    = "EXTERNAL_API_ERROR"
	CodeDatabaseError  = "DATABASE_ERROR"
	CodeTimeout        = "TIMEOUT"
)

// NotFoundError is returned when a requested resource does not exist
type NotFoundError struct {
	Resource string
	ID       interface{}
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with id '%v' not found", e.Resource, e.ID)
}

func (e *NotFoundError) Code() string {
	return CodeNotFound
}

// NewNotFoundError creates a new NotFoundError
func NewNotFoundError(resource string, id interface{}) *NotFoundError {
	return &NotFoundError{Resource: resource, ID: id}
}

// ValidationError is returned when input validation fails
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

func (e *ValidationError) Code() string {
	return CodeValidation
}

// NewValidationError creates a new ValidationError
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{Field: field, Message: message}
}

// BusinessError is returned for business logic violations
type BusinessError struct {
	ErrCode string
	Message string
}

func (e *BusinessError) Error() string {
	return e.Message
}

func (e *BusinessError) Code() string {
	return e.ErrCode
}

// NewBusinessError creates a new BusinessError
func NewBusinessError(code, message string) *BusinessError {
	return &BusinessError{ErrCode: code, Message: message}
}

// UnauthorizedError is returned when authentication is required but not provided
type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "unauthorized access"
}

func (e *UnauthorizedError) Code() string {
	return CodeUnauthorized
}

// NewUnauthorizedError creates a new UnauthorizedError
func NewUnauthorizedError(message string) *UnauthorizedError {
	return &UnauthorizedError{Message: message}
}

// ForbiddenError is returned when the user doesn't have permission
type ForbiddenError struct {
	Message string
}

func (e *ForbiddenError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "access forbidden"
}

func (e *ForbiddenError) Code() string {
	return CodeForbidden
}

// NewForbiddenError creates a new ForbiddenError
func NewForbiddenError(message string) *ForbiddenError {
	return &ForbiddenError{Message: message}
}

// ConflictError is returned when a resource already exists
type ConflictError struct {
	Resource string
	Field    string
	Value    interface{}
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("%s with %s '%v' already exists", e.Resource, e.Field, e.Value)
}

func (e *ConflictError) Code() string {
	return CodeConflict
}

// NewConflictError creates a new ConflictError
func NewConflictError(resource, field string, value interface{}) *ConflictError {
	return &ConflictError{Resource: resource, Field: field, Value: value}
}

// ExternalAPIError is returned when an external API call fails
type ExternalAPIError struct {
	Service    string
	StatusCode int
	Message    string
	Err        error
}

func (e *ExternalAPIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("external API error from %s (status %d): %s - %v", e.Service, e.StatusCode, e.Message, e.Err)
	}
	return fmt.Sprintf("external API error from %s (status %d): %s", e.Service, e.StatusCode, e.Message)
}

func (e *ExternalAPIError) Code() string {
	return CodeExternalAPI
}

func (e *ExternalAPIError) Unwrap() error {
	return e.Err
}

// NewExternalAPIError creates a new ExternalAPIError
func NewExternalAPIError(service string, statusCode int, message string, err error) *ExternalAPIError {
	return &ExternalAPIError{Service: service, StatusCode: statusCode, Message: message, Err: err}
}

// DatabaseError is returned when a database operation fails
type DatabaseError struct {
	Operation string
	Message   string
	Err       error
}

func (e *DatabaseError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("database error during %s: %s - %v", e.Operation, e.Message, e.Err)
	}
	return fmt.Sprintf("database error during %s: %s", e.Operation, e.Message)
}

func (e *DatabaseError) Code() string {
	return CodeDatabaseError
}

func (e *DatabaseError) Unwrap() error {
	return e.Err
}

// NewDatabaseError creates a new DatabaseError
func NewDatabaseError(operation, message string, err error) *DatabaseError {
	return &DatabaseError{Operation: operation, Message: message, Err: err}
}

// TimeoutError is returned when an operation times out
type TimeoutError struct {
	Operation string
	Duration  string
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("operation '%s' timed out after %s", e.Operation, e.Duration)
}

func (e *TimeoutError) Code() string {
	return CodeTimeout
}

// NewTimeoutError creates a new TimeoutError
func NewTimeoutError(operation, duration string) *TimeoutError {
	return &TimeoutError{Operation: operation, Duration: duration}
}

// CodedError interface for errors that have an error code
type CodedError interface {
	error
	Code() string
}

// GetCode returns the error code if the error implements CodedError
func GetCode(err error) string {
	var coded CodedError
	if errors.As(err, &coded) {
		return coded.Code()
	}
	return CodeInternal
}

// IsNotFound checks if the error is a NotFoundError
func IsNotFound(err error) bool {
	var notFound *NotFoundError
	return errors.As(err, &notFound)
}

// IsValidation checks if the error is a ValidationError
func IsValidation(err error) bool {
	var validation *ValidationError
	return errors.As(err, &validation)
}

// IsUnauthorized checks if the error is an UnauthorizedError
func IsUnauthorized(err error) bool {
	var unauthorized *UnauthorizedError
	return errors.As(err, &unauthorized)
}

// IsForbidden checks if the error is a ForbiddenError
func IsForbidden(err error) bool {
	var forbidden *ForbiddenError
	return errors.As(err, &forbidden)
}

// IsConflict checks if the error is a ConflictError
func IsConflict(err error) bool {
	var conflict *ConflictError
	return errors.As(err, &conflict)
}

// IsExternalAPI checks if the error is an ExternalAPIError
func IsExternalAPI(err error) bool {
	var external *ExternalAPIError
	return errors.As(err, &external)
}

// IsDatabase checks if the error is a DatabaseError
func IsDatabase(err error) bool {
	var database *DatabaseError
	return errors.As(err, &database)
}

// IsTimeout checks if the error is a TimeoutError
func IsTimeout(err error) bool {
	var timeout *TimeoutError
	return errors.As(err, &timeout)
}
