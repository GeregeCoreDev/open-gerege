// Package domain provides implementation for domain
//
// File: user.go
// Description: implementation for domain
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
/*
Package domain нь application-ийн бизнес entity-уудыг тодорхойлно.

Entity-ууд нь:
  - Database хүснэгтийн бүтэцтэй тохирно (GORM ORM)
  - JSON serialization-ийг дэмжинэ
  - Business validation агуулна

Энэ файлд User болон UserRole entity-ууд тодорхойлогдсон.

Database tables:
  - users: Хэрэглэгчийн мэдээлэл
  - user_roles: Хэрэглэгч-эрхийн холбоос (many-to-many)
*/
package domain

// ============================================================
// USER ENTITY
// ============================================================

// User нь системийн хэрэглэгчийн мэдээллийг илэрхийлнэ.
// Table: users
//
// ХУР (Хувийн мэдээллийн улсын регистр)-тэй интеграци хийх үед
// CivilId, RegNo зэрэг талбарууд ашиглагдана.
//
// GORM тайлбар:
//   - gorm:"primaryKey": Primary key
//   - gorm:"type:varchar(10)": Column type
//
// JSON тайлбар:
//   - json:"id": JSON field name
//   - json:"id,omitempty": Omit if zero value
type User struct {
	// Id нь хэрэглэгчийн primary key.
	// SSO системээс ирсэн user_id-тай тохирно.
	Id int `json:"id" gorm:"primaryKey"`

	// CivilId нь иргэний регистрийн ID.
	// ХУР системтэй холбох үед ашиглагдана.
	CivilId int `json:"civil_id"`

	// RegNo нь регистрийн дугаар (ЖШ: АА00112233).
	// Монгол иргэний регистр: 2 үсэг + 8 тоо = 10 тэмдэгт.
	RegNo string `json:"reg_no" gorm:"type:varchar(10)"`

	// FamilyName нь овог (ургийн овог).
	// Жишээ: "Борjiигин"
	FamilyName string `json:"family_name" gorm:"type:varchar(80)"`

	// LastName нь эцэг/эхийн нэр.
	// Жишээ: "Бат"
	LastName string `json:"last_name" gorm:"type:varchar(150)"`

	// FirstName нь өөрийн нэр.
	// Жишээ: "Болд"
	FirstName string `json:"first_name" gorm:"type:varchar(150)"`

	// Gender нь хүйс.
	// 1 = Эрэгтэй, 2 = Эмэгтэй
	Gender int `json:"gender"`

	// BirthDate нь төрсөн огноо (YYYY-MM-DD format).
	// Жишээ: "1990-01-15"
	BirthDate string `json:"birth_date" gorm:"type:varchar(10)"`

	// PhoneNo нь утасны дугаар (8 оронтой).
	// Жишээ: "99112233"
	PhoneNo string `json:"phone_no" gorm:"type:varchar(8)"`

	// Email нь цахим шуудан.
	// Жишээ: "user@example.com"
	Email string `json:"email" gorm:"type:varchar(80)"`

	// Status нь хэрэглэгчийн төлөв.
	// active, suspended, locked, pending_verification, deactivated
	Status string `json:"status" gorm:"default:active"`

	// StatusReason нь төлөв өөрчлөгдсөн шалтгаан
	StatusReason string `json:"status_reason"`

	// StatusChangedAt нь төлөв өөрчлөгдсөн огноо
	StatusChangedAt *LocalDateTime `json:"status_changed_at"`

	// StatusChangedBy нь төлөв өөрчилсэн хэрэглэгчийн ID
	StatusChangedBy *int `json:"status_changed_by"`

	// LastLoginAt нь сүүлд нэвтэрсэн огноо
	LastLoginAt *LocalDateTime `json:"last_login_at"`

	// LoginCount нь нийт нэвтэрсэн тоо
	LoginCount int `json:"login_count" gorm:"default:0"`

	// ExtraFields нь нийтлэг талбаруудыг агуулна:
	// - CreatedDate: Үүсгэсэн огноо
	// - UpdatedDate: Шинэчилсэн огноо
	// - CreatedBy: Үүсгэсэн хэрэглэгч
	// - UpdatedBy: Шинэчилсэн хэрэглэгч
	// - DeletedDate: Устгасан огноо (soft delete)
	// - DeletedBy: Устгасан хэрэглэгч
	ExtraFields
}

// TableName нь GORM-д хүснэгтийн нэрийг зааж өгнө.
// Default-аар "users" гэж таамаглана.
// func (User) TableName() string { return "users" }

// ============================================================
// USER ROLE ENTITY (Many-to-Many)
// ============================================================

// UserRole нь хэрэглэгч-эрхийн холбоосыг илэрхийлнэ.
// Table: user_roles
//
// Энэ нь many-to-many холбоос:
// - Нэг хэрэглэгч олон эрхтэй байж болно
// - Нэг эрх олон хэрэглэгчид хуваарилагдаж болно
//
// Composite primary key: (user_id, role_id)
//
// GORM Foreign Key тайлбар:
//   - foreignKey:UserId: Энэ struct-ийн UserId талбар
//   - references:Id: Target struct-ийн Id талбар
//   - OnUpdate:CASCADE: Target update хийхэд энэ мөр ч update хийгдэнэ
//   - OnDelete:SET NULL: Target устгахад энэ талбар NULL болно
type UserRole struct {
	// UserId нь хэрэглэгчийн ID (users.id руу FK).
	UserId int `json:"user_id"`

	// User нь холбогдсон хэрэглэгчийн мэдээлэл.
	// GORM-ийн Preload("User") ашиглаж авна.
	// json:"user,omitempty": Хоосон бол JSON-д оруулахгүй.
	User *User `json:"user,omitempty" gorm:"foreignKey:UserId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	// RoleID нь эрхийн ID (roles.id руу FK).
	RoleID int `json:"role_id"`

	// Role нь холбогдсон эрхийн мэдээлэл.
	// GORM-ийн Preload("Role") ашиглаж авна.
	Role *Role `json:"role,omitempty" gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	// ExtraFields нь нийтлэг timestamp талбаруудыг агуулна.
	ExtraFields
}

// TableName нь GORM-д хүснэгтийн нэрийг зааж өгнө.
// func (UserRole) TableName() string { return "user_roles" }
