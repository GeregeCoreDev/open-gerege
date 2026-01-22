// Package domain provides business entities
//
// File: domain_test.go
// Description: Unit tests for domain package
package domain

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalDateTime_String(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "valid date",
			time:     time.Date(2024, 1, 15, 9, 30, 0, 0, time.Local),
			expected: "2024-01-15 09:30:00",
		},
		{
			name:     "zero value",
			time:     time.Time{},
			expected: "0000-00-00 00:00:00",
		},
		{
			name:     "end of year",
			time:     time.Date(2024, 12, 31, 23, 59, 59, 0, time.Local),
			expected: "2024-12-31 23:59:59",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ldt := LocalDateTime(tt.time)
			assert.Equal(t, tt.expected, ldt.String())
		})
	}
}

func TestLocalDateTime_MarshalJSON(t *testing.T) {
	tm := time.Date(2024, 1, 15, 9, 30, 0, 0, time.Local)
	ldt := LocalDateTime(tm)

	data, err := json.Marshal(ldt)

	require.NoError(t, err)
	assert.Equal(t, `"2024-01-15 09:30:00"`, string(data))
}

func TestLocalDateTime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid date",
			input:   `"2024-01-15 09:30:00"`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ldt LocalDateTime
			err := json.Unmarshal([]byte(tt.input), &ldt)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "2024-01-15 09:30:00", ldt.String())
			}
		})
	}
}

func TestLocalDateTime_Value(t *testing.T) {
	tm := time.Date(2024, 1, 15, 9, 30, 0, 0, time.Local)
	ldt := LocalDateTime(tm)

	val, err := ldt.Value()

	require.NoError(t, err)
	assert.IsType(t, time.Time{}, val)
}

func TestLocalDateTime_Value_ZeroValue(t *testing.T) {
	ldt := LocalDateTime{}

	val, err := ldt.Value()

	require.NoError(t, err)
	// Zero value should return current time
	assert.IsType(t, time.Time{}, val)
}

func TestLocalDateTime_Scan(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "scan time.Time",
			input:    time.Date(2024, 1, 15, 9, 30, 0, 0, time.Local),
			expected: "2024-01-15 09:30:00",
		},
		{
			name:     "scan string",
			input:    "2024-01-15 09:30:00",
			expected: "2024-01-15 09:30:00",
		},
		{
			name:     "scan nil",
			input:    nil,
			expected: "0000-00-00 00:00:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ldt LocalDateTime
			err := ldt.Scan(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.expected, ldt.String())
		})
	}
}

func TestExtraFields_Structure(t *testing.T) {
	// Test that ExtraFields can be embedded
	type TestEntity struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		ExtraFields
	}

	entity := TestEntity{
		ID:   1,
		Name: "Test",
	}

	// ExtraFields should be accessible
	assert.Equal(t, 1, entity.ID)
	assert.Equal(t, "Test", entity.Name)
	assert.Nil(t, entity.CreatedDate)
	assert.Equal(t, 0, entity.CreatedUserId)
}

func TestExtraFields_JSON(t *testing.T) {
	type TestEntity struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		ExtraFields
	}

	now := LocalDateTime(time.Now())
	entity := TestEntity{
		ID:          1,
		Name:        "Test",
		ExtraFields: ExtraFields{
			CreatedDate:   &now,
			CreatedUserId: 100, // Should not appear in JSON (json:"-")
		},
	}

	data, err := json.Marshal(entity)
	require.NoError(t, err)

	// CreatedUserId should not be in JSON output
	assert.NotContains(t, string(data), "created_user_id")
	assert.Contains(t, string(data), "created_date")
}

func TestUser_Structure(t *testing.T) {
	user := User{
		Id:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}

	assert.Equal(t, 1, user.Id)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Equal(t, "john@example.com", user.Email)
}

func TestRole_Structure(t *testing.T) {
	isActive := true
	role := Role{
		ID:          1,
		Code:        "admin",
		Name:        "Administrator",
		Description: "Full access",
		IsActive:    &isActive,
	}

	assert.Equal(t, 1, role.ID)
	assert.Equal(t, "admin", role.Code)
	assert.Equal(t, "Administrator", role.Name)
	assert.True(t, *role.IsActive)
}

func TestPermission_Structure(t *testing.T) {
	perm := Permission{
		ID:   1,
		Code: "admin.user.read",
		Name: "Read Users",
	}

	assert.Equal(t, 1, perm.ID)
	assert.Equal(t, "admin.user.read", perm.Code)
	assert.Equal(t, "Read Users", perm.Name)
}

func TestOrganization_Structure(t *testing.T) {
	org := Organization{
		Id:    1,
		Name:  "Test Org",
		RegNo: "1234567",
	}

	assert.Equal(t, 1, org.Id)
	assert.Equal(t, "Test Org", org.Name)
	assert.Equal(t, "1234567", org.RegNo)
}

func TestNews_Structure(t *testing.T) {
	news := News{
		Id:       1,
		Title:    "Test Title",
		Text:     "Test content",
		ImageUrl: "https://example.com/image.jpg",
	}

	assert.Equal(t, 1, news.Id)
	assert.Equal(t, "Test Title", news.Title)
	assert.Equal(t, "Test content", news.Text)
	assert.Equal(t, "https://example.com/image.jpg", news.ImageUrl)
}

func TestModule_Structure(t *testing.T) {
	isActive := true
	module := Module{
		ID:          1,
		Code:        "user",
		Name:        "User Module",
		Description: "User management",
		IsActive:    &isActive,
		SystemID:    1,
	}

	assert.Equal(t, 1, module.ID)
	assert.Equal(t, "user", module.Code)
	assert.Equal(t, "User Module", module.Name)
	assert.True(t, *module.IsActive)
}

func TestSystem_Structure(t *testing.T) {
	isActive := true
	system := System{
		ID:          1,
		Code:        "admin",
		Name:        "Admin System",
		Description: "Administration",
		IsActive:    &isActive,
	}

	assert.Equal(t, 1, system.ID)
	assert.Equal(t, "admin", system.Code)
	assert.Equal(t, "Admin System", system.Name)
	assert.True(t, *system.IsActive)
}
