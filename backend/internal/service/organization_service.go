// Package service provides implementation for service
//
// File: organization_service.go
// Description: implementation for service
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package service

import (
	"context"
	"fmt"
	"net/url"

	"templatev25/internal/domain"
	"templatev25/internal/http/dto"
	"templatev25/internal/repository"

	"git.gerege.mn/backend-packages/common"
	"git.gerege.mn/backend-packages/config"
	"git.gerege.mn/backend-packages/httpx"
	"git.gerege.mn/backend-packages/utils"
	"go.uber.org/zap"
)

type OrganizationService struct {
	repo repository.OrganizationRepository
	log  *zap.Logger
}

func NewOrganizationService(repo repository.OrganizationRepository, log *zap.Logger) *OrganizationService {
	return &OrganizationService{repo: repo, log: log}
}

func (s *OrganizationService) List(ctx context.Context, p common.PaginationQuery) ([]domain.Organization, int64, int, int, error) {
	items, total, page, size, err := s.repo.List(ctx, p)
	if err != nil {
		s.log.Error("organization_list_failed", zap.Error(err))
		return nil, 0, 0, 0, err
	}
	s.log.Debug("organization_list_success", zap.Int64("total", total), zap.Int("page", page))
	return items, total, page, size, nil
}

func (s *OrganizationService) Create(ctx context.Context, req dto.OrganizationDto) (domain.Organization, error) {
	// defaults
	if req.ShortName == "" {
		req.ShortName = req.Name
	}
	if req.ParentID != nil && *req.ParentID == 0 {
		req.ParentID = nil
	}
	m := domain.Organization{
		Id:                req.Id,
		RegNo:             req.RegNo,
		Name:              req.Name,
		ShortName:         req.ShortName,
		TypeId:            req.TypeId,
		PhoneNo:           req.PhoneNo,
		Email:             req.Email,
		Longitude:         req.Longitude,
		Latitude:          req.Latitude,
		IsActive:          req.IsActive,
		AimagId:           req.AimagId,
		SumId:             req.SumId,
		BagId:             req.BagId,
		AddressDetail:     req.AddressDetail,
		AimagName:         req.AimagName,
		SumName:           req.SumName,
		BagName:           req.BagName,
		CountryCode:       req.CountryCode,
		CountryName:       req.CountryName,
		Sequence:          req.Sequence,
		ParentAddressId:   req.ParentAddressId,
		ParentAddressName: req.ParentAddressName,
		CountryNameEn:     req.CountryNameEn,
		ParentId:          req.ParentID,
	}
	org, err := s.repo.Create(ctx, m)
	if err != nil {
		s.log.Error("organization_create_failed", zap.String("name", req.Name), zap.Error(err))
		return domain.Organization{}, err
	}
	s.log.Info("organization_created", zap.Int("org_id", org.Id), zap.String("name", org.Name))
	return org, nil
}

func (s *OrganizationService) Update(ctx context.Context, id int, req dto.OrganizationUpdateDto) (domain.Organization, error) {
	if req.ShortName == "" {
		req.ShortName = req.Name
	}
	if req.ParentID != nil && *req.ParentID == 0 {
		req.ParentID = nil
	}
	m := domain.Organization{
		RegNo:             req.RegNo,
		Name:              req.Name,
		ShortName:         req.ShortName,
		TypeId:            req.TypeId,
		PhoneNo:           req.PhoneNo,
		Email:             req.Email,
		Longitude:         req.Longitude,
		Latitude:          req.Latitude,
		IsActive:          req.IsActive,
		AimagId:           req.AimagId,
		SumId:             req.SumId,
		BagId:             req.BagId,
		AddressDetail:     req.AddressDetail,
		AimagName:         req.AimagName,
		SumName:           req.SumName,
		BagName:           req.BagName,
		CountryCode:       req.CountryCode,
		CountryName:       req.CountryName,
		Sequence:          req.Sequence,
		ParentAddressId:   req.ParentAddressId,
		ParentAddressName: req.ParentAddressName,
		CountryNameEn:     req.CountryNameEn,
		ParentId:          req.ParentID,
	}
	org, err := s.repo.Update(ctx, id, m)
	if err != nil {
		s.log.Error("organization_update_failed", zap.Int("org_id", id), zap.Error(err))
		return domain.Organization{}, err
	}
	s.log.Info("organization_updated", zap.Int("org_id", id))
	return org, nil
}

func (s *OrganizationService) Delete(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.log.Error("organization_delete_failed", zap.Int("org_id", id), zap.Error(err))
		return err
	}
	s.log.Info("organization_deleted", zap.Int("org_id", id))
	return nil
}

func (s *OrganizationService) ByID(ctx context.Context, id int) (domain.Organization, error) {
	org, err := s.repo.ByID(ctx, id)
	if err != nil {
		s.log.Error("organization_get_by_id_failed", zap.Int("org_id", id), zap.Error(err))
		return domain.Organization{}, err
	}
	return org, nil
}

func (s *OrganizationService) Tree(ctx context.Context, rootID int) ([]domain.Organization, error) {
	items, err := s.repo.Tree(ctx, rootID)
	if err != nil {
		s.log.Error("organization_tree_failed", zap.Int("root_id", rootID), zap.Error(err))
		return nil, err
	}
	s.log.Debug("organization_tree_fetched", zap.Int("root_id", rootID), zap.Int("count", len(items)))
	return items, nil
}

type OrganizationTypeService struct {
	repo repository.OrganizationTypeRepository
}

func NewOrganizationTypeService(repo repository.OrganizationTypeRepository) *OrganizationTypeService {
	return &OrganizationTypeService{repo: repo}
}

func (s *OrganizationTypeService) List(ctx context.Context, p common.PaginationQuery) ([]domain.OrganizationType, int64, int, int, error) {
	return s.repo.List(ctx, p)
}

func (s *OrganizationTypeService) Create(ctx context.Context, req dto.OrganizationTypeDto) error {
	m := domain.OrganizationType{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
	}
	return s.repo.Create(ctx, m)
}

func (s *OrganizationTypeService) Update(ctx context.Context, id int, req dto.OrganizationTypeDto) error {
	m := domain.OrganizationType{Code: req.Code, Name: req.Name, Description: req.Description}
	return s.repo.Update(ctx, id, m)
}

func (s *OrganizationTypeService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *OrganizationTypeService) Systems(ctx context.Context, typeID int) ([]domain.System, error) {
	return s.repo.Systems(ctx, typeID)
}

func (s *OrganizationTypeService) AddSystems(ctx context.Context, typeID int, systemIDs []int) error {
	return s.repo.AddSystems(ctx, typeID, systemIDs)
}

func (s *OrganizationTypeService) Roles(ctx context.Context, typeID int) ([]domain.Role, error) {
	return s.repo.Roles(ctx, typeID)
}

func (s *OrganizationTypeService) AddRoles(ctx context.Context, typeID int, roleIDs []int) error {
	return s.repo.AddRoles(ctx, typeID, roleIDs)
}

type OrgUserService struct {
	repo  repository.OrgUserRepository
	urepo repository.UserRepository
	http  *httpx.Client
	cfg   *config.Config
}

func NewOrgUserService(repo repository.OrgUserRepository, cfg *config.Config, urepo repository.UserRepository) *OrgUserService {
	return &OrgUserService{
		repo:  repo,
		cfg:   cfg,
		http:  httpx.New(0),
		urepo: urepo,
	}
}

func (s *OrgUserService) List(ctx context.Context, q dto.OrgUserListQuery) ([]domain.OrganizationUser, int64, int, int, error) {
	return s.repo.List(ctx, q)
}

func (s *OrgUserService) Add(ctx context.Context, req dto.OrgUserCreateDto, authHeader string) error {
	// давхардал шалгах
	if _, err := s.repo.FindByOrgAndUser(ctx, req.OrgId, req.UserId); err == nil {
		return fmt.Errorf("хэрэглэгч аль хэдийн бүртгэгдсэн байна")
	}

	// байгууллага байх ёстой
	ok, err := s.repo.OrgExists(ctx, req.OrgId)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("байгууллага олдсонгүй")
	}

	// хэрэглэгч локалд байхгүй бол CORE → citizen/find
	uok, err := s.repo.UserExists(ctx, req.UserId)
	if err != nil {
		return err
	}
	if !uok {
		endpoint := fmt.Sprintf("%s/citizen/find?search_text=%s", s.cfg.URLS.Core, url.QueryEscape(fmt.Sprintf("%d", req.UserId)))
		// Core талын хариуг ашиглан локал DB-д бүртгэдэг өөр модуль/handler танайд байгаа тул энд зөвхөн fetch-ийг гүйцэтгэнэ.
		var _ any
		resp, _, err := httpx.GetJSON[dto.CoreUser](ctx, s.http, endpoint, map[string]string{
			"Authorization": authHeader,
		})

		if err != nil {
			return fmt.Errorf("хэрэглэгчийн мэдээлэл олдсонгүй")
		}

		m := domain.User{
			Id:         resp.Id,
			CivilId:    resp.CivilId,
			RegNo:      resp.RegNo,
			FamilyName: resp.FamilyName,
			LastName:   resp.LastName,
			FirstName:  resp.FirstName,
			Gender:     resp.Gender,
			BirthDate:  utils.NormalizeDate(resp.BirthDate),
			PhoneNo:    resp.PhoneNo,
			Email:      resp.Email,
		}
		// exists check (хуучин логик)
		if _, err := s.urepo.GetByID(ctx, resp.Id); err != nil {
			if _, err := s.urepo.Create(ctx, m); err != nil {
				return err
			}
		}
		// Тайлбар: Хэрэв локал insert шаардлагатай бол энд User insert хийх логикоо нэмээрэй.
	}

	return s.repo.Add(ctx, domain.OrganizationUser{
		OrgId:  req.OrgId,
		UserId: req.UserId,
	})
}

func (s *OrgUserService) Remove(ctx context.Context, req dto.OrgUserDeleteDto) error {
	return s.repo.Remove(ctx, req.OrgId, req.UserId)
}

func (s *OrgUserService) UsersByOrg(ctx context.Context, orgId int, name string, p common.PaginationQuery) ([]dto.ResOrguserUserItem, int64, int, int, error) {
	page, size, offset := utils.OffsetLimit(p) // reuse helper
	_ = offset
	items, total, err := s.repo.ListUsersByOrg(ctx, orgId, name, page, size)
	return items, total, page, size, err
}

func (s *OrgUserService) OrgsByUser(ctx context.Context, userId int, name string, p common.PaginationQuery) ([]dto.ResOrguserOrgItem, int64, int, int, error) {
	page, size, offset := utils.OffsetLimit(p)
	_ = offset
	items, total, err := s.repo.ListOrgsByUser(ctx, userId, name, page, size)
	return items, total, page, size, err
}
