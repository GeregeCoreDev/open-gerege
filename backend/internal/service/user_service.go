// Package service provides implementation for service
//
// File: user_service.go
// Description: implementation for service
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package service

import (
	"context"

	"templatev25/internal/domain"
	"templatev25/internal/http/dto"
	"templatev25/internal/middleware"
	"templatev25/internal/repository"

	"git.gerege.mn/backend-packages/common"
	"git.gerege.mn/backend-packages/config"
	"git.gerege.mn/backend-packages/utils"
	"go.uber.org/zap"
)

type UserService struct {
	repo repository.UserRepository
	log  *zap.Logger
	cfg  *config.Config
}

func NewUserService(repo repository.UserRepository, cfg *config.Config, log *zap.Logger) *UserService {
	return &UserService{
		repo: repo,
		log:  log,
		cfg:  cfg,
	}
}

func (s *UserService) GetByID(ctx context.Context, id int) (domain.User, error) {
	log := middleware.LoggerOrDefault(ctx, s.log)
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Error("user_get_by_id_failed", zap.Int("user_id", id), zap.Error(err))
		return domain.User{}, err
	}
	return user, nil
}

// List — PaginationQuery (model_repo хэв маяг)
func (s *UserService) List(ctx context.Context, p common.PaginationQuery) ([]domain.User, int64, int, int, error) {
	log := middleware.LoggerOrDefault(ctx, s.log)
	items, total, page, size, err := s.repo.List(ctx, p)
	if err != nil {
		log.Error("user_list_failed", zap.Error(err))
		return nil, 0, 0, 0, err
	}
	log.Debug("user_list_success", zap.Int64("total", total), zap.Int("page", page))
	return items, total, page, size, nil
}

func (s *UserService) Create(ctx context.Context, req dto.UserCreateDto) (domain.User, error) {
	log := middleware.LoggerOrDefault(ctx, s.log)
	m := domain.User{
		Id:         req.Id,
		CivilId:    req.CivilId,
		RegNo:      req.RegNo,
		FamilyName: req.FamilyName,
		LastName:   req.LastName,
		FirstName:  req.FirstName,
		Gender:     req.Gender,
		BirthDate:  utils.NormalizeDate(req.BirthDate),
		PhoneNo:    req.PhoneNo,
		Email:      req.Email,
	}
	// exists check (хуучин логик)
	if user, err := s.repo.GetByID(ctx, req.Id); err == nil {
		log.Debug("user_already_exists", zap.Int("user_id", req.Id))
		return user, nil
	}
	user, err := s.repo.Create(ctx, m)
	if err != nil {
		log.Error("user_create_failed", zap.Int("user_id", req.Id), zap.Error(err))
		return domain.User{}, err
	}
	log.Info("user_created", zap.Int("user_id", user.Id), zap.String("reg_no", user.RegNo))
	return user, nil
}

func (s *UserService) Update(ctx context.Context, req dto.UserUpdateDto) (domain.User, error) {
	log := middleware.LoggerOrDefault(ctx, s.log)
	// exists check
	if _, err := s.repo.GetByID(ctx, req.Id); err != nil {
		log.Error("user_update_not_found", zap.Int("user_id", req.Id), zap.Error(err))
		return domain.User{}, err
	}
	m := domain.User{
		Id:         req.Id,
		CivilId:    req.CivilId,
		RegNo:      req.RegNo,
		FamilyName: req.FamilyName,
		LastName:   req.LastName,
		FirstName:  req.FirstName,
		Gender:     req.Gender,
		BirthDate:  utils.NormalizeDate(req.BirthDate),
		PhoneNo:    req.PhoneNo,
		Email:      req.Email,
	}
	user, err := s.repo.Update(ctx, m)
	if err != nil {
		log.Error("user_update_failed", zap.Int("user_id", req.Id), zap.Error(err))
		return domain.User{}, err
	}
	log.Info("user_updated", zap.Int("user_id", user.Id))
	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id int) (domain.User, error) {
	log := middleware.LoggerOrDefault(ctx, s.log)
	user, err := s.repo.Delete(ctx, id)
	if err != nil {
		log.Error("user_delete_failed", zap.Int("user_id", id), zap.Error(err))
		return domain.User{}, err
	}
	log.Info("user_deleted", zap.Int("user_id", id))
	return user, nil
}

// -------- Profile & Organizations --------

func (s *UserService) Organizations(ctx context.Context, userID, currentOrgID int, fields []string) (orgID int, org *domain.Organization, items []domain.Organization, err error) {
	log := middleware.LoggerOrDefault(ctx, s.log)
	ids, err := s.repo.UserOrgIDs(ctx, userID)
	if err != nil {
		log.Error("user_orgs_ids_failed", zap.Int("user_id", userID), zap.Error(err))
		return 0, nil, nil, err
	}
	items, err = s.repo.GetOrganizationsByIDs(ctx, ids, fields)
	if err != nil {
		log.Error("user_orgs_get_failed", zap.Int("user_id", userID), zap.Error(err))
		return 0, nil, nil, err
	}
	if currentOrgID > 0 {
		org, err = s.repo.GetOrganization(ctx, currentOrgID, fields)
		if err != nil {
			log.Error("user_current_org_failed", zap.Int("org_id", currentOrgID), zap.Error(err))
			return 0, nil, nil, err
		}
		orgID = currentOrgID
	}
	log.Debug("user_orgs_fetched", zap.Int("user_id", userID), zap.Int("org_count", len(items)))
	return orgID, org, items, nil
}
