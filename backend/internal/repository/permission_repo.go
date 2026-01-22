// Package repository provides implementation for repository
//
// File: permission_repo.go
// Description: implementation for repository
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package repository

import (
	"context"
	"strings"
	"templatev25/internal/domain"
	"templatev25/internal/http/dto"

	"time"

	"git.gerege.mn/backend-packages/ctx"
	"git.gerege.mn/backend-packages/scopes"
	"git.gerege.mn/backend-packages/utils"

	"gorm.io/gorm"
)

type PermissionRepository interface {
	List(ctx context.Context, q dto.PermissionQuery) ([]domain.Permission, int64, int, int, error)
	ByID(ctx context.Context, id int) (domain.Permission, error)
	ByCode(ctx context.Context, code string) (domain.Permission, error)
	Create(ctx context.Context, m domain.Permission) error
	CreateBatch(ctx context.Context, systemID int, moduleID int, actionIDs []int64) error
	Update(ctx context.Context, id int, m domain.Permission) error
	Delete(ctx context.Context, id int) error

	// Permission шалгах методууд
	UserHasPermission(ctx context.Context, userID int, permissionCode string) (bool, error)
	GetUserPermissionCodes(ctx context.Context, userID int) ([]string, error)
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) List(ctx context.Context, q dto.PermissionQuery) ([]domain.Permission, int64, int, int, error) {
	page, size, offset := utils.OffsetLimit(q.PaginationQuery)

	colMap := scopes.ColumnMap{
		"id":          "permissions.id",
		"code":        "permissions.code",
		"name":        "permissions.name",
		"description": "permissions.description",
		"module_id":   "permissions.module_id",
		"system_id":   "permissions.system_id",
		"action_id":   "permissions.action_id",
	}

	tx := r.db.WithContext(ctx).Model(&domain.Permission{}).Scopes(
		scopes.SearchScope(colMap, utils.ParseSearch(q.Search)),
		scopes.DateScope(q.CreatedFrom, q.CreatedTo),
	)

	if q.SystemID > 0 {
		tx = tx.Where("system_id = ?", q.SystemID)
	}

	if q.ModuleID > 0 {
		tx = tx.Where("module_id = ?", q.ModuleID)
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, 0, 0, err
	}

	var items []domain.Permission
	if err := tx.Scopes(
		scopes.SortScope(colMap, utils.ParseSort(q.Sort), "id DESC"),
	).Preload("Module").Preload("System").Preload("Action").Offset(offset).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, 0, 0, err
	}
	return items, total, page, size, nil
}

func (r *permissionRepository) ByID(ctx context.Context, id int) (domain.Permission, error) {
	var m domain.Permission
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return domain.Permission{}, err
	}
	return m, nil
}

func (r *permissionRepository) ByCode(ctx context.Context, code string) (domain.Permission, error) {
	var m domain.Permission
	if err := r.db.WithContext(ctx).Where("code = ? AND deleted_date IS NULL", code).First(&m).Error; err != nil {
		return domain.Permission{}, err
	}
	return m, nil
}

func (r *permissionRepository) Create(uctx context.Context, m domain.Permission) error {
	if uid, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.CreatedUserId = uid
	}
	if oid, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.CreatedOrgId = oid
	}
	return r.db.WithContext(uctx).Create(&m).Error
}

func (r *permissionRepository) CreateBatch(uctx context.Context, systemID int, moduleID int, actionIDs []int64) error {
	// ctx-оос CreatedUser/Org онооно
	var createdUserId, createdOrgId int
	if uid, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		createdUserId = uid
	}
	if oid, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		createdOrgId = oid
	}

	// Transaction ашиглаж бүх Permission-г нэгэн зэрэг үүсгэх
	return WithTx(uctx, r.db, func(tx *gorm.DB) error {
		// System-ийн code-г олох
		var system domain.System
		if err := tx.Where("id = ?", systemID).First(&system).Error; err != nil {
			return err
		}

		// Module-ийн code-г олох
		var module domain.Module
		if err := tx.Where("id = ?", moduleID).First(&module).Error; err != nil {
			return err
		}

		// Action-уудын мэдээллийг авах
		var actions []domain.Action
		if err := tx.Where("id IN ?", actionIDs).Find(&actions).Error; err != nil {
			return err
		}

		// Action-уудын тоо шалгах
		if len(actions) != len(actionIDs) {
			return gorm.ErrRecordNotFound
		}

		// Permission-ууд үүсгэх (Action бүрт нэг Permission)
		for _, action := range actions {
			// Permission code-г systemcode.modulecode.actioncode гэж үүсгэх (lower case)
			permissionCode := strings.ToLower(system.Code) + "." + strings.ToLower(module.Code) + "." + strings.ToLower(action.Code)

			permission := domain.Permission{
				Code:        permissionCode,
				Name:        action.Name,
				Description: action.Description,
				SystemID:    systemID,
				ModuleID:    moduleID,
				ActionID:    &action.ID,
				IsActive:    action.IsActive,
			}

			// Permission-ийн CreatedUser/Org-ийг тохируулах
			permission.CreatedUserId = createdUserId
			permission.CreatedOrgId = createdOrgId

			// Permission үүсгэх
			if err := tx.Create(&permission).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *permissionRepository) Update(uctx context.Context, id int, m domain.Permission) error {
	if uid, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.UpdatedUserId = uid
	}
	if oid, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.UpdatedOrgId = oid
	}
	return r.db.WithContext(uctx).Model(&domain.Permission{}).Where("id = ?", id).Updates(&m).Error
}

func (r *permissionRepository) Delete(uctx context.Context, id int) error {
	m := domain.Permission{}
	if uid, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.DeletedUserId = uid
	}
	if oid, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.DeletedOrgId = oid
	}
	m.DeletedDate = gorm.DeletedAt{Valid: true, Time: time.Now()}
	return r.db.WithContext(uctx).Where("id = ?", id).Updates(&m).Error
}

// UserHasPermission нь хэрэглэгч тодорхой permission-тэй эсэхийг шалгана.
// user_roles -> roles -> role_permissions -> permissions гэсэн холбоосоор шалгана.
//
// Parameters:
//   - ctx: Context
//   - userID: Хэрэглэгчийн ID
//   - permissionCode: Permission код (жишээ: "admin.role.create")
//
// Returns:
//   - bool: Permission байвал true
//   - error: Алдаа
func (r *permissionRepository) UserHasPermission(ctx context.Context, userID int, permissionCode string) (bool, error) {
	var exists bool
	err := r.db.WithContext(ctx).Raw(`
		SELECT EXISTS(
			SELECT 1 FROM permissions p
			JOIN role_permissions rp ON p.id = rp.permission_id
			JOIN user_roles ur ON ur.role_id = rp.role_id
			WHERE ur.user_id = ?
			AND p.code = ?
			AND p.is_active = true
			AND p.deleted_date IS NULL
			AND rp.deleted_date IS NULL
			AND ur.deleted_date IS NULL
		)
	`, userID, permissionCode).Scan(&exists).Error
	if err != nil {
		return false, err
	}
	return exists, nil
}

// GetUserPermissionCodes нь хэрэглэгчийн бүх permission код-уудыг буцаана.
// user_roles -> roles -> role_permissions -> permissions гэсэн холбоосоор авна.
//
// Parameters:
//   - ctx: Context
//   - userID: Хэрэглэгчийн ID
//
// Returns:
//   - []string: Permission кодуудын жагсаалт
//   - error: Алдаа
func (r *permissionRepository) GetUserPermissionCodes(ctx context.Context, userID int) ([]string, error) {
	var codes []string
	err := r.db.WithContext(ctx).Raw(`
		SELECT DISTINCT p.code FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN user_roles ur ON ur.role_id = rp.role_id
		WHERE ur.user_id = ?
		AND p.is_active = true
		AND p.deleted_date IS NULL
		AND rp.deleted_date IS NULL
		AND ur.deleted_date IS NULL
	`, userID).Scan(&codes).Error
	if err != nil {
		return nil, err
	}
	return codes, nil
}
