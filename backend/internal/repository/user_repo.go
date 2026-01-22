// Package repository provides implementation for repository
//
// File: user_repo.go
// Description: implementation for repository
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package repository

import (
	"context"
	"time"

	"templatev25/internal/domain"

	"git.gerege.mn/backend-packages/common"
	"git.gerege.mn/backend-packages/ctx"
	"git.gerege.mn/backend-packages/scopes"
	"git.gerege.mn/backend-packages/utils"

	"gorm.io/gorm"
)

type UserRepository interface {
	List(ctx context.Context, p common.PaginationQuery) ([]domain.User, int64, int, int, error)
	Create(ctx context.Context, m domain.User) (domain.User, error)
	Update(ctx context.Context, m domain.User) (domain.User, error)
	Delete(ctx context.Context, id int) (domain.User, error)
	GetByID(ctx context.Context, id int) (domain.User, error)

	// Organizations helper (profile/organizations endpoint-д хэрэглэнэ)
	UserOrgIDs(ctx context.Context, userID int) ([]int, error)
	GetOrganizationsByIDs(ctx context.Context, ids []int, fields []string) ([]domain.Organization, error)
	GetOrganization(ctx context.Context, id int, fields []string) (*domain.Organization, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// List — model_repo хэв маяг: scopes + pagination + олон талбарт name хайлт
func (r *userRepository) List(ctx context.Context, p common.PaginationQuery) ([]domain.User, int64, int, int, error) {
	page, size, offset := utils.OffsetLimit(p)

	colMap := scopes.ColumnMap{
		"id":          "users.id",
		"reg_no":      "users.reg_no",
		"first_name":  "users.first_name",
		"last_name":   "users.last_name",
		"phone_no":    "users.phone_no",
		"email":       "users.email",
		"birth_date":  "users.birth_date",
		"civil_id":    "users.civil_id",
		"family_name": "users.family_name",
		"gender":      "users.gender",
	}

	tx := r.db.WithContext(ctx).Model(&domain.User{}).Scopes(
		scopes.SearchScope(colMap, utils.ParseSearch(p.Search)),
		scopes.DateScope(p.CreatedFrom, p.CreatedTo),
	)

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, 0, 0, err
	}

	var items []domain.User
	if err := tx.Scopes(
		// Хуучин “name” хайлтыг орлуулахын тулд first/last/phone/reg талбаруудыг default-д оруулсан
		scopes.SortScope(colMap, utils.ParseSort(p.Sort), "id DESC"),
	).Offset(offset).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, 0, 0, err
	}

	return items, total, page, size, nil
}

func (r *userRepository) Create(ctx context.Context, m domain.User) (domain.User, error) {
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return domain.User{}, err
	}
	return m, nil
}

func (r *userRepository) Update(ctx context.Context, m domain.User) (domain.User, error) {

	if err := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", m.Id).
		Updates(&m).Error; err != nil {
		return domain.User{}, err
	}
	return m, nil
}

func (r *userRepository) Delete(uctx context.Context, id int) (domain.User, error) {
	// Soft delete - RoleRepository-тай адил pattern
	var ex domain.User
	if err := r.db.WithContext(uctx).Take(&ex, "id = ?", id).Error; err != nil {
		return domain.User{}, err
	}

	// DeletedUser/Org context-оос авах
	m := domain.User{}
	if userId, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.DeletedUserId = userId
	}
	if orgId, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.DeletedOrgId = orgId
	}
	m.DeletedDate = gorm.DeletedAt{Valid: true, Time: time.Now()}

	if err := r.db.WithContext(uctx).
		Model(&domain.User{}).
		Where("id = ?", id).
		Updates(&m).Error; err != nil {
		return domain.User{}, err
	}

	return ex, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int) (domain.User, error) {
	var u domain.User
	err := r.db.WithContext(ctx).Take(&u, "id = ?", id).Error
	return u, err
}

// ---------- Organizations helpers ----------

func (r *userRepository) UserOrgIDs(ctx context.Context, userID int) ([]int, error) {
	var ids []int
	err := r.db.WithContext(ctx).Model(&domain.OrganizationUser{}).
		Where("user_id = ?", userID).
		Pluck("org_id", &ids).Error
	return ids, err
}

func (r *userRepository) GetOrganizationsByIDs(ctx context.Context, ids []int, fields []string) ([]domain.Organization, error) {
	var out []domain.Organization
	tx := r.db.WithContext(ctx).Model(&domain.Organization{})
	if len(fields) > 0 {
		tx = tx.Select(fields)
	}
	err := tx.Where("id in (?)", ids).Find(&out).Error
	return out, err
}

func (r *userRepository) GetOrganization(ctx context.Context, id int, fields []string) (*domain.Organization, error) {
	var o domain.Organization
	tx := r.db.WithContext(ctx).Model(&domain.Organization{})
	if len(fields) > 0 {
		tx = tx.Select(fields)
	}
	if err := tx.Take(&o, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &o, nil
}
