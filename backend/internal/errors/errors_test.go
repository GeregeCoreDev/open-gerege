// Package errors provides custom error types
//
// File: errors_test.go
// Description: Unit tests for custom error types
package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotFoundError(t *testing.T) {
	err := NewNotFoundError("User", 123)

	assert.Equal(t, "User with id '123' not found", err.Error())
	assert.Equal(t, CodeNotFound, err.Code())
	assert.Equal(t, "User", err.Resource)
	assert.Equal(t, 123, err.ID)
}

func TestNotFoundError_StringID(t *testing.T) {
	err := NewNotFoundError("Article", "abc-123")

	assert.Equal(t, "Article with id 'abc-123' not found", err.Error())
}

func TestValidationError(t *testing.T) {
	err := NewValidationError("email", "invalid email format")

	assert.Equal(t, "validation error on field 'email': invalid email format", err.Error())
	assert.Equal(t, CodeValidation, err.Code())
	assert.Equal(t, "email", err.Field)
	assert.Equal(t, "invalid email format", err.Message)
}

func TestValidationError_NoField(t *testing.T) {
	err := NewValidationError("", "general validation error")

	assert.Equal(t, "validation error: general validation error", err.Error())
}

func TestBusinessError(t *testing.T) {
	err := NewBusinessError("INSUFFICIENT_BALANCE", "not enough funds")

	assert.Equal(t, "not enough funds", err.Error())
	assert.Equal(t, "INSUFFICIENT_BALANCE", err.Code())
}

func TestUnauthorizedError(t *testing.T) {
	err := NewUnauthorizedError("invalid token")

	assert.Equal(t, "invalid token", err.Error())
	assert.Equal(t, CodeUnauthorized, err.Code())
}

func TestUnauthorizedError_DefaultMessage(t *testing.T) {
	err := NewUnauthorizedError("")

	assert.Equal(t, "unauthorized access", err.Error())
}

func TestForbiddenError(t *testing.T) {
	err := NewForbiddenError("insufficient permissions")

	assert.Equal(t, "insufficient permissions", err.Error())
	assert.Equal(t, CodeForbidden, err.Code())
}

func TestForbiddenError_DefaultMessage(t *testing.T) {
	err := NewForbiddenError("")

	assert.Equal(t, "access forbidden", err.Error())
}

func TestConflictError(t *testing.T) {
	err := NewConflictError("User", "email", "test@example.com")

	assert.Equal(t, "User with email 'test@example.com' already exists", err.Error())
	assert.Equal(t, CodeConflict, err.Code())
	assert.Equal(t, "User", err.Resource)
	assert.Equal(t, "email", err.Field)
	assert.Equal(t, "test@example.com", err.Value)
}

func TestExternalAPIError(t *testing.T) {
	innerErr := errors.New("connection refused")
	err := NewExternalAPIError("PaymentService", 500, "payment failed", innerErr)

	assert.Contains(t, err.Error(), "external API error from PaymentService")
	assert.Contains(t, err.Error(), "status 500")
	assert.Contains(t, err.Error(), "payment failed")
	assert.Contains(t, err.Error(), "connection refused")
	assert.Equal(t, CodeExternalAPI, err.Code())
	assert.Equal(t, innerErr, err.Unwrap())
}

func TestExternalAPIError_NoInnerError(t *testing.T) {
	err := NewExternalAPIError("AuthService", 401, "unauthorized", nil)

	assert.Equal(t, "external API error from AuthService (status 401): unauthorized", err.Error())
	assert.Nil(t, err.Unwrap())
}

func TestDatabaseError(t *testing.T) {
	innerErr := errors.New("connection timeout")
	err := NewDatabaseError("insert", "failed to insert user", innerErr)

	assert.Contains(t, err.Error(), "database error during insert")
	assert.Contains(t, err.Error(), "failed to insert user")
	assert.Contains(t, err.Error(), "connection timeout")
	assert.Equal(t, CodeDatabaseError, err.Code())
	assert.Equal(t, innerErr, err.Unwrap())
}

func TestDatabaseError_NoInnerError(t *testing.T) {
	err := NewDatabaseError("query", "invalid query", nil)

	assert.Equal(t, "database error during query: invalid query", err.Error())
	assert.Nil(t, err.Unwrap())
}

func TestTimeoutError(t *testing.T) {
	err := NewTimeoutError("API call", "30s")

	assert.Equal(t, "operation 'API call' timed out after 30s", err.Error())
	assert.Equal(t, CodeTimeout, err.Code())
}

func TestGetCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{"NotFoundError", NewNotFoundError("Test", 1), CodeNotFound},
		{"ValidationError", NewValidationError("field", "msg"), CodeValidation},
		{"UnauthorizedError", NewUnauthorizedError("msg"), CodeUnauthorized},
		{"ForbiddenError", NewForbiddenError("msg"), CodeForbidden},
		{"ConflictError", NewConflictError("R", "F", "V"), CodeConflict},
		{"ExternalAPIError", NewExternalAPIError("S", 500, "msg", nil), CodeExternalAPI},
		{"DatabaseError", NewDatabaseError("op", "msg", nil), CodeDatabaseError},
		{"TimeoutError", NewTimeoutError("op", "10s"), CodeTimeout},
		{"BusinessError", NewBusinessError("CUSTOM", "msg"), "CUSTOM"},
		{"StandardError", errors.New("standard error"), CodeInternal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, GetCode(tt.err))
		})
	}
}

func TestIsNotFound(t *testing.T) {
	assert.True(t, IsNotFound(NewNotFoundError("Test", 1)))
	assert.False(t, IsNotFound(NewValidationError("field", "msg")))
	assert.False(t, IsNotFound(errors.New("standard error")))
}

func TestIsValidation(t *testing.T) {
	assert.True(t, IsValidation(NewValidationError("field", "msg")))
	assert.False(t, IsValidation(NewNotFoundError("Test", 1)))
	assert.False(t, IsValidation(errors.New("standard error")))
}

func TestIsUnauthorized(t *testing.T) {
	assert.True(t, IsUnauthorized(NewUnauthorizedError("msg")))
	assert.False(t, IsUnauthorized(NewForbiddenError("msg")))
	assert.False(t, IsUnauthorized(errors.New("standard error")))
}

func TestIsForbidden(t *testing.T) {
	assert.True(t, IsForbidden(NewForbiddenError("msg")))
	assert.False(t, IsForbidden(NewUnauthorizedError("msg")))
	assert.False(t, IsForbidden(errors.New("standard error")))
}

func TestIsConflict(t *testing.T) {
	assert.True(t, IsConflict(NewConflictError("R", "F", "V")))
	assert.False(t, IsConflict(NewNotFoundError("Test", 1)))
	assert.False(t, IsConflict(errors.New("standard error")))
}

func TestIsExternalAPI(t *testing.T) {
	assert.True(t, IsExternalAPI(NewExternalAPIError("S", 500, "msg", nil)))
	assert.False(t, IsExternalAPI(NewDatabaseError("op", "msg", nil)))
	assert.False(t, IsExternalAPI(errors.New("standard error")))
}

func TestIsDatabase(t *testing.T) {
	assert.True(t, IsDatabase(NewDatabaseError("op", "msg", nil)))
	assert.False(t, IsDatabase(NewExternalAPIError("S", 500, "msg", nil)))
	assert.False(t, IsDatabase(errors.New("standard error")))
}

func TestIsTimeout(t *testing.T) {
	assert.True(t, IsTimeout(NewTimeoutError("op", "10s")))
	assert.False(t, IsTimeout(NewDatabaseError("op", "msg", nil)))
	assert.False(t, IsTimeout(errors.New("standard error")))
}

func TestErrorCodes(t *testing.T) {
	// Verify all error codes are defined
	assert.Equal(t, "NOT_FOUND", CodeNotFound)
	assert.Equal(t, "VALIDATION_ERROR", CodeValidation)
	assert.Equal(t, "UNAUTHORIZED", CodeUnauthorized)
	assert.Equal(t, "FORBIDDEN", CodeForbidden)
	assert.Equal(t, "CONFLICT", CodeConflict)
	assert.Equal(t, "INTERNAL_ERROR", CodeInternal)
	assert.Equal(t, "BAD_REQUEST", CodeBadRequest)
	assert.Equal(t, "EXTERNAL_API_ERROR", CodeExternalAPI)
	assert.Equal(t, "DATABASE_ERROR", CodeDatabaseError)
	assert.Equal(t, "TIMEOUT", CodeTimeout)
}
