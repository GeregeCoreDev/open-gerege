// Package repository provides implementation for repository
//
// File: organization_repo.go
// Description: implementation for repository
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"templatev25/internal/domain"
	"templatev25/internal/http/dto"

	"git.gerege.mn/backend-packages/common"
	"git.gerege.mn/backend-packages/config"

	"git.gerege.mn/backend-packages/ctx"
	"git.gerege.mn/backend-packages/scopes"
	"git.gerege.mn/backend-packages/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrganizationRepository interface {
	List(ctx context.Context, p common.PaginationQuery) ([]domain.Organization, int64, int, int, error)
	Create(ctx context.Context, m domain.Organization) (domain.Organization, error)
	Update(ctx context.Context, id int, m domain.Organization) (domain.Organization, error)
	Delete(ctx context.Context, id int) error
	ByID(ctx context.Context, id int) (domain.Organization, error)
	Tree(ctx context.Context, rootID int) ([]domain.Organization, error)
}

type organizationRepository struct{ db *gorm.DB }

func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) List(ctx context.Context, p common.PaginationQuery) ([]domain.Organization, int64, int, int, error) {
	page, size, offset := utils.OffsetLimit(p)
	colMap := scopes.ColumnMap{
		"id":         "organizations.id",
		"name":       "organizations.name",
		"short_name": "organizations.short_name",
		"reg_no":     "organizations.reg_no",
		"type_id":    "organizations.type_id",
	}
	tx := r.db.WithContext(ctx).Model(&domain.Organization{}).
		Preload("Type").
		Scopes(scopes.SearchScope(colMap, utils.ParseSearch(p.Search)),
			scopes.DateScope(p.CreatedFrom, p.CreatedTo))

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, 0, 0, err
	}

	var items []domain.Organization
	if err := tx.Scopes(scopes.SortScope(colMap, utils.ParseSort(p.Sort), "name ASC")).
		Offset(offset).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, 0, 0, err
	}
	return items, total, page, size, nil
}

func (r *organizationRepository) Create(ctx context.Context, m domain.Organization) (domain.Organization, error) {
	if err := r.db.WithContext(ctx).Clauses(clause.Returning{}, clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&m).Error; err != nil {
		return domain.Organization{}, err
	}
	return m, nil
}

func (r *organizationRepository) Update(ctx context.Context, id int, m domain.Organization) (domain.Organization, error) {
	m.Id = id
	if err := r.db.WithContext(ctx).Clauses(clause.Returning{}).
		Model(&domain.Organization{}).
		Where("id = ?", id).
		Updates(&m).Error; err != nil {
		return domain.Organization{}, err
	}
	return m, nil
}

func (r *organizationRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&domain.OrganizationUser{}, "org_id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Delete(&domain.Organization{Id: id}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *organizationRepository) ByID(ctx context.Context, id int) (domain.Organization, error) {
	var o domain.Organization
	err := r.db.WithContext(ctx).Preload("Type").Take(&o, "id = ?", id).Error
	return o, err
}

func (r *organizationRepository) Tree(ctx context.Context, rootID int) ([]domain.Organization, error) {
	var items []domain.Organization
	// Хэрэв танайд ParentPreloader/ChildrenPreloader байгаа бол түүнийг хэрэглээрэй.
	if err := r.db.WithContext(ctx).
		Preload("Children").
		Find(&items, "id = ?", rootID).Error; err != nil {
		return nil, err
	}
	return items, nil
}

type OrganizationTypeRepository interface {
	List(ctx context.Context, p common.PaginationQuery) ([]domain.OrganizationType, int64, int, int, error)
	Create(ctx context.Context, m domain.OrganizationType) error
	Update(ctx context.Context, id int, m domain.OrganizationType) error
	Delete(ctx context.Context, id int) error

	// System linkage
	AddSystems(ctx context.Context, orgTypeID int, systemIDs []int) error
	Systems(ctx context.Context, orgTypeID int) ([]domain.System, error)

	// Role linkage
	AddRoles(ctx context.Context, orgTypeID int, roleIDs []int) error
	Roles(ctx context.Context, orgTypeID int) ([]domain.Role, error)
}

type organizationTypeRepository struct{ db *gorm.DB }

func NewOrganizationTypeRepository(db *gorm.DB) OrganizationTypeRepository {
	return &organizationTypeRepository{db: db}
}

func (r *organizationTypeRepository) List(ctx context.Context, p common.PaginationQuery) ([]domain.OrganizationType, int64, int, int, error) {
	page, size, offset := utils.OffsetLimit(p)
	colMap := scopes.ColumnMap{
		"id":   "organization_types.id",
		"name": "organization_types.name",
	}
	tx := r.db.WithContext(ctx).Model(&domain.OrganizationType{}).Scopes(
		scopes.SearchScope(colMap, utils.ParseSearch(p.Search)),
		scopes.DateScope(p.CreatedFrom, p.CreatedTo),
	)

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, 0, 0, err
	}
	var items []domain.OrganizationType
	if err := tx.Scopes(scopes.SortScope(colMap, utils.ParseSort(p.Sort), "id DESC")).
		Offset(offset).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, 0, 0, err
	}
	return items, total, page, size, nil
}

func (r *organizationTypeRepository) Create(uctx context.Context, m domain.OrganizationType) error {
	if userId, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.CreatedUserId = userId
	}
	if orgId, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.CreatedOrgId = orgId
	}

	return r.db.WithContext(uctx).Create(&m).Error
}

func (r *organizationTypeRepository) Update(uctx context.Context, id int, m domain.OrganizationType) error {
	if userId, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.UpdatedUserId = userId
	}
	if orgId, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.UpdatedOrgId = orgId
	}
	return r.db.WithContext(uctx).
		Model(&domain.OrganizationType{}).
		Where("id = ?", id).
		Updates(&m).Error
}

func (r *organizationTypeRepository) Delete(uctx context.Context, id int) error {
	m := domain.OrganizationType{}
	if userId, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.DeletedUserId = userId
	}
	if orgId, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.DeletedOrgId = orgId
	}
	m.DeletedDate = gorm.DeletedAt{Valid: true, Time: time.Now()}
	return r.db.WithContext(uctx).Where("id = ?", id).Updates(&m).Error
}

type OrgUserRepository interface {
	// generic list (org_id or user_id-р шүүнэ, name filter нь тухайн preload дээр хамаарна)
	List(ctx context.Context, q dto.OrgUserListQuery) ([]domain.OrganizationUser, int64, int, int, error)

	// joins
	ListUsersByOrg(ctx context.Context, orgId int, name string, page, size int) ([]dto.ResOrguserUserItem, int64, error)
	ListOrgsByUser(ctx context.Context, userId int, name string, page, size int) ([]dto.ResOrguserOrgItem, int64, error)

	Add(ctx context.Context, ou domain.OrganizationUser) error
	Remove(ctx context.Context, orgId, userId int) error

	OrgExists(ctx context.Context, orgId int) (bool, error)
	UserExists(ctx context.Context, userId int) (bool, error)
	FindByOrgAndUser(ctx context.Context, orgId, userId int) (domain.OrganizationUser, error)
}

type orgUserRepository struct {
	db *gorm.DB
}

func NewOrgUserRepository(db *gorm.DB, cfg *config.Config) OrgUserRepository {
	// cfg parameter kept for backward compatibility but no longer used
	// search_path is now set in DSN, so schema name is not needed
	return &orgUserRepository{db: db}
}

func (r *orgUserRepository) List(ctx context.Context, q dto.OrgUserListQuery) ([]domain.OrganizationUser, int64, int, int, error) {
	page, size, offset := utils.OffsetLimit(q.PaginationQuery)

	var (
		items    []domain.OrganizationUser
		total    int64
		nameLike = strings.TrimSpace(q.Name)
	)

	cnt := r.db.WithContext(ctx).Model(&domain.OrganizationUser{})
	tx := r.db.WithContext(ctx).Model(&domain.OrganizationUser{}).Order("created_date DESC")

	if q.UserId != 0 {
		cnt = cnt.Where("user_id = ?", q.UserId)
		tx = tx.Where("user_id = ?", q.UserId).
			Preload("Organization", func(db *gorm.DB) *gorm.DB {
				if nameLike != "" {
					n := "%" + nameLike + "%"
					db = db.Where("name ILIKE ? OR short_name ILIKE ? OR reg_no ILIKE ?", n, n, n)
				}
				return db
			})
	}
	if q.OrgId != 0 {
		cnt = cnt.Where("org_id = ?", q.OrgId)
		tx = tx.Where("org_id = ?", q.OrgId).
			Preload("User", func(db *gorm.DB) *gorm.DB {
				if nameLike != "" {
					n := "%" + nameLike + "%"
					db = db.Where("first_name ILIKE ? OR last_name ILIKE ? OR reg_no ILIKE ? OR phone_no ILIKE ?", n, n, n, n)
				}
				return db
			})
	}

	if err := cnt.Count(&total).Error; err != nil {
		return nil, 0, 0, 0, err
	}
	if err := tx.Offset(offset).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, 0, 0, err
	}
	return items, total, page, size, nil
}

func (r *orgUserRepository) Add(ctx context.Context, ou domain.OrganizationUser) error {
	return r.db.WithContext(ctx).Create(&ou).Error
}

func (r *orgUserRepository) Remove(ctx context.Context, orgId, userId int) error {
	return r.db.WithContext(ctx).Delete(&domain.OrganizationUser{}, "org_id = ? AND user_id = ?", orgId, userId).Error
}

func (r *orgUserRepository) OrgExists(ctx context.Context, orgId int) (bool, error) {
	var cnt int64
	if err := r.db.WithContext(ctx).Model(&domain.Organization{}).Where("id = ?", orgId).Count(&cnt).Error; err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (r *orgUserRepository) UserExists(ctx context.Context, userId int) (bool, error) {
	var cnt int64
	if err := r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", userId).Count(&cnt).Error; err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (r *orgUserRepository) FindByOrgAndUser(ctx context.Context, orgId, userId int) (domain.OrganizationUser, error) {
	var m domain.OrganizationUser
	err := r.db.WithContext(ctx).Where("org_id = ? AND user_id = ?", orgId, userId).First(&m).Error
	return m, err
}

// ---------- Raw JOIN queries (pagination гарыг нь удирдана) ----------

func (r *orgUserRepository) ListUsersByOrg(ctx context.Context, orgId int, name string, page, size int) ([]dto.ResOrguserUserItem, int64, error) {
	var (
		total int64
		rows  []dto.ResOrguserUserItem
	)

	nameCond := ""
	argsCnt := []any{orgId}
	if strings.TrimSpace(name) != "" {
		nameCond = "AND (tu.first_name ILIKE @name OR tu.last_name ILIKE @name OR tu.reg_no ILIKE @name OR tu.phone_no ILIKE @name)"
		argsCnt = append(argsCnt, sql.Named("name", "%"+name+"%"))
	}

	cntSQL := fmt.Sprintf(`
		SELECT COUNT(*) FROM organization_users tou
		LEFT JOIN users tu ON tou.user_id = tu.id
		WHERE tou.deleted_date IS NULL
		  AND tu.deleted_date IS NULL
		  AND tou.org_id = ?
		  %s
	`, nameCond)
	if err := r.db.WithContext(ctx).Raw(cntSQL, argsCnt...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	querySQL := fmt.Sprintf(`
		SELECT
			tou.org_id,
			tu.id AS user_id,
			tu.last_name,
			tu.first_name,
			tu.reg_no,
			tu.birth_date,
			tu.gender,
			tu.phone_no,
			tu.email,
			tou.created_date
		FROM organization_users tou
		LEFT JOIN users tu ON tou.user_id = tu.id
		WHERE tou.deleted_date IS NULL
		  AND tu.deleted_date IS NULL
		  AND tou.org_id = ?
		  %s
		LIMIT %d OFFSET %d
	`, nameCond, size, offset)

	args := []any{orgId}
	if strings.TrimSpace(name) != "" {
		args = append(args, sql.Named("name", "%"+name+"%"))
	}
	if err := r.db.WithContext(ctx).Raw(querySQL, args...).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *orgUserRepository) ListOrgsByUser(ctx context.Context, userId int, name string, page, size int) ([]dto.ResOrguserOrgItem, int64, error) {
	var (
		total int64
		rows  []dto.ResOrguserOrgItem
	)

	nameCond := ""
	argsCnt := []any{userId}
	if strings.TrimSpace(name) != "" {
		nameCond = "AND (tu.name ILIKE @name OR tu.short_name ILIKE @name OR tu.reg_no ILIKE @name)"
		argsCnt = append(argsCnt, sql.Named("name", "%"+name+"%"))
	}

	cntSQL := fmt.Sprintf(`
		SELECT COUNT(*) FROM organization_users tou
		LEFT JOIN organizations tu ON tou.org_id = tu.id
		WHERE tou.deleted_date IS NULL
		  AND tu.deleted_date IS NULL
		  AND tou.user_id = ?
		  %s
	`, nameCond)
	if err := r.db.WithContext(ctx).Raw(cntSQL, argsCnt...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	querySQL := fmt.Sprintf(`
		SELECT
			tou.org_id,
			tu.id,
			tu.name,
			tu.short_name,
			tu.reg_no,
			tou.created_date
		FROM organization_users tou
		LEFT JOIN organizations tu ON tou.org_id = tu.id
		WHERE tou.deleted_date IS NULL
		  AND tu.deleted_date IS NULL
		  AND tou.user_id = ?
		  %s
		LIMIT %d OFFSET %d
	`, nameCond, size, offset)

	args := []any{userId}
	if strings.TrimSpace(name) != "" {
		args = append(args, sql.Named("name", "%"+name+"%"))
	}
	if err := r.db.WithContext(ctx).Raw(querySQL, args...).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// --- System linkage ---

func (r *organizationTypeRepository) Systems(ctx context.Context, orgTypeID int) ([]domain.System, error) {
	var links []domain.OrgTypeSystem
	if err := r.db.WithContext(ctx).
		Preload("System").
		Where("type_id = ?", orgTypeID).
		Find(&links).Error; err != nil {
		return nil, err
	}

	out := make([]domain.System, 0, len(links))
	for _, l := range links {
		if l.System != nil {
			out = append(out, *l.System)
		}
	}
	return out, nil
}

func (r *organizationTypeRepository) AddSystems(ctx context.Context, orgTypeID int, systemIDs []int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// одоогийн map-уудыг цэвэрлээд шинээр үүсгэнэ (replace semantics)
		if err := tx.WithContext(ctx).
			Where("type_id = ?", orgTypeID).
			Delete(&domain.OrgTypeSystem{}).Error; err != nil {
			return err
		}

		if len(systemIDs) == 0 {
			return nil
		}

		// Batch insert instead of N individual inserts
		links := make([]domain.OrgTypeSystem, 0, len(systemIDs))
		for _, sid := range systemIDs {
			links = append(links, domain.OrgTypeSystem{
				TypeId:   orgTypeID,
				SystemID: sid,
			})
		}
		return tx.WithContext(ctx).Create(&links).Error
	})
}

// --- Role linkage ---

func (r *organizationTypeRepository) Roles(ctx context.Context, orgTypeID int) ([]domain.Role, error) {
	var links []domain.OrgTypeRole
	if err := r.db.WithContext(ctx).
		Preload("Role").
		Where("type_id = ?", orgTypeID).
		Find(&links).Error; err != nil {
		return nil, err
	}

	out := make([]domain.Role, 0, len(links))
	for _, l := range links {
		if l.Role != nil {
			out = append(out, *l.Role)
		}
	}
	return out, nil
}

func (r *organizationTypeRepository) AddRoles(ctx context.Context, orgTypeID int, roleIDs []int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// одоогийн map-уудыг цэвэрлээд шинээр үүсгэнэ (replace semantics)
		if err := tx.WithContext(ctx).
			Where("type_id = ?", orgTypeID).
			Delete(&domain.OrgTypeRole{}).Error; err != nil {
			return err
		}

		if len(roleIDs) == 0 {
			return nil
		}

		// Batch insert instead of N individual inserts
		links := make([]domain.OrgTypeRole, 0, len(roleIDs))
		for _, rid := range roleIDs {
			links = append(links, domain.OrgTypeRole{
				TypeId: orgTypeID,
				RoleID: rid,
			})
		}
		return tx.WithContext(ctx).Create(&links).Error
	})
}
