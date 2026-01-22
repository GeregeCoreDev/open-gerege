// Package errors provides custom error types for the application
//
// File: errors_test.go
// Description: Unit tests for custom error types
package errors_test

import (
	"errors"
	"fmt"
	"testing"

	apperrors "templatev25/internal/errors"

	"github.com/stretchr/testify/assert"
)

func TestNotFoundError(t *testing.T) {
	tests := []struct {
		name     string
		resource string
		id       interface{}
		wantMsg  string
	}{
		{
			name:     "with int id",
			resource: "User",
			id:       1,
			wantMsg:  "User with id '1' not found",
		},
		{
			name:     "with string id",
			resource: "Role",
			id:       "admin",
			wantMsg:  "Role with id 'admin' not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := apperrors.NewNotFoundError(tt.resource, tt.id)

			assert.Equal(t, tt.wantMsg, err.Error())
			assert.Equal(t, apperrors.CodeNotFound, err.Code())
			assert.True(t, apperrors.IsNotFound(err))
		})
	}
}

func TestValidationError(t *testing.T) {
	tests := []struct {
		name    string
		field   string
		message string
		wantMsg string
	}{
		{
			name:    "with field",
			field:   "email",
			message: "invalid format",
			wantMsg: "validation error on field 'email': invalid format",
		},
		{
			name:    "without field",
			field:   "",
			message: "invalid input",
			wantMsg: "validation error: invalid input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := apperrors.NewValidationError(tt.field, tt.message)

			assert.Equal(t, tt.wantMsg, err.Error())
			assert.Equal(t, apperrors.CodeValidation, err.Code())
			assert.True(t, apperrors.IsValidation(err))
		})
	}
}

func TestBusinessError(t *testing.T) {
	err := apperrors.NewBusinessError(apperrors.CodeBadRequest, "operation not allowed")

	assert.Equal(t, "operation not allowed", err.Error())
	assert.Equal(t, apperrors.CodeBadRequest, err.Code())
}

func TestUnauthorizedError(t *testing.T) {
	tests := []struct {
		name    string
		message string
		wantMsg string
	}{
		{
			name:    "with custom message",
			message: "token expired",
			wantMsg: "token expired",
		},
		{
			name:    "without message",
			message: "",
			wantMsg: "unauthorized access",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := apperrors.NewUnauthorizedError(tt.message)

			assert.Equal(t, tt.wantMsg, err.Error())
			assert.Equal(t, apperrors.CodeUnauthorized, err.Code())
			assert.True(t, apperrors.IsUnauthorized(err))
		})
	}
}

func TestForbiddenError(t *testing.T) {
	tests := []struct {
		name    string
		message string
		wantMsg string
	}{
		{
			name:    "with custom message",
			message: "insufficient permissions",
			wantMsg: "insufficient permissions",
		},
		{
			name:    "without message",
			message: "",
			wantMsg: "access forbidden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := apperrors.NewForbiddenError(tt.message)

			assert.Equal(t, tt.wantMsg, err.Error())
			assert.Equal(t, apperrors.CodeForbidden, err.Code())
			assert.True(t, apperrors.IsForbidden(err))
		})
	}
}

func TestConflictError(t *testing.T) {
	err := apperrors.NewConflictError("User", "email", "test@example.com")

	assert.Equal(t, "User with email 'test@example.com' already exists", err.Error())
	assert.Equal(t, apperrors.CodeConflict, err.Code())
	assert.True(t, apperrors.IsConflict(err))
}

func TestExternalAPIError(t *testing.T) {
	tests := []struct {
		name       string
		service    string
		statusCode int
		message    string
		err        error
		wantMsg    string
	}{
		{
			name:       "with underlying error",
			service:    "PaymentAPI",
			statusCode: 500,
			message:    "server error",
			err:        fmt.Errorf("connection refused"),
			wantMsg:    "external API error from PaymentAPI (status 500): server error - connection refused",
		},
		{
			name:       "without underlying error",
			service:    "UserAPI",
			statusCode: 404,
			message:    "user not found",
			err:        nil,
			wantMsg:    "external API error from UserAPI (status 404): user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := apperrors.NewExternalAPIError(tt.service, tt.statusCode, tt.message, tt.err)

			assert.Equal(t, tt.wantMsg, err.Error())
			assert.Equal(t, apperrors.CodeExternalAPI, err.Code())
			assert.True(t, apperrors.IsExternalAPI(err))

			if tt.err != nil {
				assert.Equal(t, tt.err, err.Unwrap())
			}
		})
	}
}

func TestDatabaseError(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		message   string
		err       error
		wantMsg   string
	}{
		{
			name:      "with underlying error",
			operation: "INSERT",
			message:   "duplicate key",
			err:       fmt.Errorf("unique constraint violation"),
			wantMsg:   "database error during INSERT: duplicate key - unique constraint violation",
		},
		{
			name:      "without underlying error",
			operation: "SELECT",
			message:   "connection failed",
			err:       nil,
			wantMsg:   "database error during SELECT: connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := apperrors.NewDatabaseError(tt.operation, tt.message, tt.err)

			assert.Equal(t, tt.wantMsg, err.Error())
			assert.Equal(t, apperrors.CodeDatabaseError, err.Code())
			assert.True(t, apperrors.IsDatabase(err))

			if tt.err != nil {
				assert.Equal(t, tt.err, err.Unwrap())
			}
		})
	}
}

func TestTimeoutError(t *testing.T) {
	err := apperrors.NewTimeoutError("API call", "30s")

	assert.Equal(t, "operation 'API call' timed out after 30s", err.Error())
	assert.Equal(t, apperrors.CodeTimeout, err.Code())
	assert.True(t, apperrors.IsTimeout(err))
}

func TestGetCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode string
	}{
		{
			name:     "not found error",
			err:      apperrors.NewNotFoundError("User", 1),
			wantCode: apperrors.CodeNotFound,
		},
		{
			name:     "validation error",
			err:      apperrors.NewValidationError("email", "invalid"),
			wantCode: apperrors.CodeValidation,
		},
		{
			name:     "regular error",
			err:      errors.New("some error"),
			wantCode: apperrors.CodeInternal,
		},
		{
			name:     "wrapped coded error",
			err:      fmt.Errorf("wrapped: %w", apperrors.NewNotFoundError("User", 1)),
			wantCode: apperrors.CodeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := apperrors.GetCode(tt.err)
			assert.Equal(t, tt.wantCode, code)
		})
	}
}

func TestIsChecks(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		checkFunc func(error) bool
		want      bool
	}{
		{
			name:      "IsNotFound - true",
			err:       apperrors.NewNotFoundError("User", 1),
			checkFunc: apperrors.IsNotFound,
			want:      true,
		},
		{
			name:      "IsNotFound - false",
			err:       apperrors.NewValidationError("field", "msg"),
			checkFunc: apperrors.IsNotFound,
			want:      false,
		},
		{
			name:      "IsValidation - true",
			err:       apperrors.NewValidationError("field", "msg"),
			checkFunc: apperrors.IsValidation,
			want:      true,
		},
		{
			name:      "IsUnauthorized - true",
			err:       apperrors.NewUnauthorizedError(""),
			checkFunc: apperrors.IsUnauthorized,
			want:      true,
		},
		{
			name:      "IsForbidden - true",
			err:       apperrors.NewForbiddenError(""),
			checkFunc: apperrors.IsForbidden,
			want:      true,
		},
		{
			name:      "IsConflict - true",
			err:       apperrors.NewConflictError("User", "email", "test@test.com"),
			checkFunc: apperrors.IsConflict,
			want:      true,
		},
		{
			name:      "IsExternalAPI - true",
			err:       apperrors.NewExternalAPIError("API", 500, "error", nil),
			checkFunc: apperrors.IsExternalAPI,
			want:      true,
		},
		{
			name:      "IsDatabase - true",
			err:       apperrors.NewDatabaseError("INSERT", "error", nil),
			checkFunc: apperrors.IsDatabase,
			want:      true,
		},
		{
			name:      "IsTimeout - true",
			err:       apperrors.NewTimeoutError("operation", "10s"),
			checkFunc: apperrors.IsTimeout,
			want:      true,
		},
		{
			name:      "regular error - all false",
			err:       errors.New("regular error"),
			checkFunc: apperrors.IsNotFound,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.checkFunc(tt.err)
			assert.Equal(t, tt.want, result)
		})
	}
}
