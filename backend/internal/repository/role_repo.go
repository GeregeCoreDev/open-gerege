// Package repository provides implementation for repository
//
// File: role_repo.go
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
	"templatev25/internal/http/dto"

	"git.gerege.mn/backend-packages/ctx"
	"git.gerege.mn/backend-packages/scopes"
	"git.gerege.mn/backend-packages/utils"

	"gorm.io/gorm"
)

type RoleRepository interface {
	// model_repo шиг PaginationQuery дамжуулдаг
	List(ctx context.Context, p dto.RoleListQuery) ([]domain.Role, int64, int, int, error)
	ByID(ctx context.Context, id int) (domain.Role, error)
	// model_repo-ийн signature-тэй тааруулсан
	Create(ctx context.Context, m domain.Role) error
	Update(ctx context.Context, id int, m domain.Role) error
	Delete(ctx context.Context, id int) error

	Permissions(ctx context.Context, q dto.RolePermissionsQuery) ([]domain.Permission, error)
	ReplacePermissions(ctx context.Context, roleID int, permIDs []int) error
	GetUserCount(uctx context.Context, id int) int64
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) ByID(ctx context.Context, id int) (domain.Role, error) {
	var m domain.Role
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return domain.Role{}, err
	}
	return m, nil
}

// -----------------------------------------------------------------------------
// List — model_repo List-тэй ижил structure (scopes + pagination)
// -----------------------------------------------------------------------------
func (r *roleRepository) List(uctx context.Context, p dto.RoleListQuery) ([]domain.Role, int64, int, int, error) {
	page, size, offset := utils.OffsetLimit(p.PaginationQuery)

	colMap := scopes.ColumnMap{
		"id":          "roles.id",
		"name":        "roles.name",
		"description": "roles.description",
	}

	tx := r.db.WithContext(uctx).Model(&domain.Role{}).Scopes(
		scopes.SearchScope(colMap, utils.ParseSearch(p.Search)),
		scopes.DateScope(p.CreatedFrom, p.CreatedTo),
	)

	if p.SystemId > 0 {
		tx = tx.Where("system_id = ?", p.SystemId)
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, 0, 0, err
	}

	var items []domain.Role
	if err := tx.Scopes(
		scopes.SortScope(colMap, utils.ParseSort(p.Sort), "id DESC"),
	).Offset(offset).Limit(size).Preload("System").Find(&items).Error; err != nil {
		return nil, 0, 0, 0, err
	}

	return items, total, page, size, nil
}

// -----------------------------------------------------------------------------
// Create/Update/Delete — model_repo-ийн convention-ийг дагана
// -----------------------------------------------------------------------------

func (r *roleRepository) Create(uctx context.Context, m domain.Role) error {
	// CreatedUser/Org context-оос
	if userId, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.CreatedUserId = userId
	}
	if orgId, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.CreatedOrgId = orgId
	}

	return r.db.WithContext(uctx).Create(&m).Error
}

func (r *roleRepository) Update(uctx context.Context, id int, m domain.Role) error {
	// UpdatedUser/Org context-оос
	if userId, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.UpdatedUserId = userId
	}
	if orgId, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.UpdatedOrgId = orgId
	}

	return r.db.WithContext(uctx).
		Model(&domain.Role{}).
		Where("id = ?", id).
		Updates(&m).Error
}

func (r *roleRepository) Delete(uctx context.Context, id int) error {
	// model_repo шиг soft-delete (suffix НЭМЭХГҮЙ!)
	m := domain.Role{}
	if userId, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.DeletedUserId = userId
	}
	if orgId, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.DeletedOrgId = orgId
	}
	m.DeletedDate = gorm.DeletedAt{Valid: true, Time: time.Now()}

	return r.db.WithContext(uctx).
		Where("id = ?", id).
		Updates(&m).Error
}

func (r *roleRepository) Permissions(ctx context.Context, q dto.RolePermissionsQuery) ([]domain.Permission, error) {
	var links []domain.RolePermission
	if err := r.db.WithContext(ctx).
		Preload("Permission.Module").
		Where("role_id = ?", q.RoleID).
		Find(&links).Error; err != nil {
		return nil, err
	}

	out := make([]domain.Permission, 0, len(links))
	for _, l := range links {
		if l.Permission != nil {
			out = append(out, *l.Permission)
		}
	}
	return out, nil
}

func (r *roleRepository) ReplacePermissions(ctx context.Context, roleID int, permIDs []int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// clear old
		if err := tx.Unscoped().Where("role_id = ?", roleID).Delete(&domain.RolePermission{}).Error; err != nil {
			return err
		}

		if len(permIDs) == 0 {
			return nil
		}

		links := make([]domain.RolePermission, 0, len(permIDs))
		for _, pid := range permIDs {
			links = append(links, domain.RolePermission{
				RoleID:       roleID,
				PermissionID: pid,
			})
		}

		if err := tx.Create(&links).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *roleRepository) GetUserCount(uctx context.Context, id int) int64 {
	cnt := int64(0)
	r.db.WithContext(uctx).Model(&domain.UserRole{}).Where("role_id = ?", id).Count(&cnt)
	return cnt
}
