// Package domain provides implementation for domain
//
// File: helper.go
// Description: implementation for domain
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
/*
Package domain нь application-ийн бизнес entity-уудыг тодорхойлно.

Энэ файл нь бүх entity-уудад нийтлэг ашиглагдах helper struct-уудыг агуулна:
  - ExtraFields: Audit талбарууд (created, updated, deleted)
  - LocalDateTime: Монгол цагийн бүс дээрх огноо

Ашиглалт:

	type User struct {
	    Id   int    `json:"id" gorm:"primaryKey"`
	    Name string `json:"name"`
	    ExtraFields  // Embed хийнэ
	}
*/
package domain

import (
	"database/sql/driver" // Database driver value interface
	"time"                // Time operations

	"gorm.io/gorm" // ORM (soft delete support)
)

// ============================================================
// EXTRA FIELDS (Audit Fields)
// ============================================================

// ExtraFields нь бүх entity-уудад нийтлэг audit талбаруудыг агуулна.
// Энэ struct-ийг embed хийснээр тухайн entity нь автоматаар
// created, updated, deleted талбаруудыг авна.
//
// GORM автомат талбарууд:
//   - autoCreateTime: INSERT хийхэд автоматаар тохируулагдана
//   - autoUpdateTime: UPDATE хийхэд автоматаар тохируулагдана
//
// Soft delete:
//   - DeletedDate нь gorm.DeletedAt төрөлтэй
//   - GORM-ийн Delete() ажиллуулахад энэ талбар тохируулагдана
//   - Find() query-ууд автоматаар deleted мөрүүдийг алгасна
//   - Unscoped() ашиглавал deleted мөрүүдийг ч авна
//
// JSON сериализаци:
//   - json:"-": JSON-д оруулахгүй (internal талбарууд)
//   - json:",omitempty": Хоосон бол JSON-д оруулахгүй
type ExtraFields struct {
	// CreatedDate нь мөр үүссэн огноо.
	// GORM autoCreateTime: INSERT хийхэд автомат тохируулна.
	CreatedDate *LocalDateTime `json:"created_date,omitempty" gorm:"autoCreateTime"`

	// CreatedUserId нь мөр үүсгэсэн хэрэглэгчийн ID.
	// Handler-ээс context-оос авч тохируулна.
	CreatedUserId int `json:"-"`

	// CreatedOrgId нь мөр үүсгэсэн байгууллагын ID.
	CreatedOrgId int `json:"-"`

	// UpdatedDate нь мөр шинэчлэгдсэн огноо.
	// GORM autoUpdateTime: UPDATE хийхэд автомат тохируулна.
	UpdatedDate *LocalDateTime `json:"updated_date,omitempty" gorm:"autoUpdateTime"`

	// UpdatedUserId нь мөр шинэчилсэн хэрэглэгчийн ID.
	UpdatedUserId int `json:"-"`

	// UpdatedOrgId нь мөр шинэчилсэн байгууллагын ID.
	UpdatedOrgId int `json:"-"`

	// DeletedUserId нь мөр устгасан хэрэглэгчийн ID.
	DeletedUserId int `json:"-"`

	// DeletedOrgId нь мөр устгасан байгууллагын ID.
	DeletedOrgId int `json:"-"`

	// DeletedDate нь soft delete огноо.
	// gorm.DeletedAt төрөл нь GORM-ийн soft delete feature-ийг идэвхжүүлнэ.
	// - Delete() → энэ талбарыг тохируулна (мөр устгахгүй)
	// - Find() → deleted_date IS NULL шүүлтүүр нэмнэ
	// - Unscoped().Delete() → жинхэнэ устгах (hard delete)
	DeletedDate gorm.DeletedAt `json:"-" gorm:"column:deleted_date;index"`
}

// ============================================================
// LOCAL DATETIME TYPE
// ============================================================

// LocalDateTime нь Монгол цагийн бүс дээрх огноог илэрхийлнэ.
// time.Time-ийн wrapper бөгөөд JSON сериализаци болон database
// хадгалалтыг тусгай format-аар хийнэ.
//
// Format: "2006-01-02 15:04:05" (Монгол стандарт)
// UTC биш, local time zone ашиглана.
type LocalDateTime time.Time

// localTimeZoneFormat нь огноо цагийн format.
// Go-ийн reference time: "Mon Jan 2 15:04:05 MST 2006"
// Бид: "2006-01-02 15:04:05" (YYYY-MM-DD HH:MM:SS)
const localTimeZoneFormat = "2006-01-02 15:04:05"

// ============================================================
// JSON MARSHALING
// ============================================================

// UnmarshalJSON нь JSON string-ийг LocalDateTime руу хөрвүүлнэ.
// Input: "2024-01-15 09:30:00"
// Output: LocalDateTime{time.Time}
//
// Implements: json.Unmarshaler interface
func (t *LocalDateTime) UnmarshalJSON(data []byte) (err error) {
	// JSON string-ийн эргэн тойрон дахь хашилтыг format-д оруулна
	now, err := time.ParseInLocation(`"`+localTimeZoneFormat+`"`, string(data), time.Local)
	*t = LocalDateTime(now)
	return
}

// MarshalJSON нь LocalDateTime-ийг JSON string руу хөрвүүлнэ.
// Input: LocalDateTime{time.Time}
// Output: "2024-01-15 09:30:00"
//
// Implements: json.Marshaler interface
func (t LocalDateTime) MarshalJSON() ([]byte, error) {
	// Pre-allocate buffer (format length + 2 quotes)
	b := make([]byte, 0, len(localTimeZoneFormat)+2)
	b = append(b, '"')
	b = append(b, []byte(t.String())...)
	b = append(b, '"')

	return b, nil
}

// ============================================================
// STRING REPRESENTATION
// ============================================================

// String нь LocalDateTime-ийг string руу хөрвүүлнэ.
// Zero value бол "0000-00-00 00:00:00" буцаана.
//
// Implements: fmt.Stringer interface
func (t LocalDateTime) String() string {
	// Zero value check
	if time.Time(t).IsZero() {
		return "0000-00-00 00:00:00"
	}

	return time.Time(t).Format(localTimeZoneFormat)
}

// ============================================================
// DATABASE VALUE (driver.Valuer interface)
// ============================================================

// Value нь LocalDateTime-ийг database-д хадгалах утга руу хөрвүүлнэ.
// Zero value бол одоогийн цаг буцаана (default value).
//
// Implements: driver.Valuer interface
func (t LocalDateTime) Value() (driver.Value, error) {
	if time.Time(t).IsZero() {
		return time.Now().Local(), nil
	}
	return time.Time(t).Local(), nil
}

// ============================================================
// DATABASE SCAN (sql.Scanner interface)
// ============================================================

// Scan нь database-ээс уншсан утгыг LocalDateTime руу хөрвүүлнэ.
// time.Time болон string төрлийг дэмжинэ.
//
// Implements: sql.Scanner interface
func (t *LocalDateTime) Scan(v interface{}) error {
	switch vt := v.(type) {
	case time.Time:
		// Database-ээс time.Time ирсэн бол шууд хөрвүүлнэ
		*t = LocalDateTime(vt)
	case string:
		// String ирсэн бол parse хийнэ
		tTime, _ := time.Parse("2006-01-02 15:04:05", vt)
		*t = LocalDateTime(tTime)
	default:
		// Бусад төрөл бол юу ч хийхгүй
		return nil
	}
	return nil
}
