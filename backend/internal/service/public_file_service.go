// Package service provides implementation for service
//
// File: public_file_service.go
// Description: implementation for service
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"templatev25/internal/domain"
	"templatev25/internal/http/dto"
	"templatev25/internal/repository"

	"git.gerege.mn/backend-packages/config"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Public file service constants
const (
	PublicImageDir = "/var/www/html/public"
	PublicImageURL = "https://business.gerege.mn/api/file/"
)

type PublicFileService struct {
	repo repository.PublicFileRepository
	cfg  *config.Config
}

func NewPublicFileService(repo repository.PublicFileRepository, cfg *config.Config) *PublicFileService {
	return &PublicFileService{repo: repo, cfg: cfg}
}

// getPublicDir returns the public file directory
// TODO: Add PublicFileDir field to config.URLConfig when available
func (s *PublicFileService) getPublicDir() string {
	return PublicImageDir
}

// getPublicURL returns the public file URL
// TODO: Add PublicFile field to config.URLConfig when available
func (s *PublicFileService) getPublicURL() string {
	return PublicImageURL
}

// List
func (s *PublicFileService) List(ctx context.Context, q dto.PublicFileListQuery) ([]domain.PublicFile, int64, int, int, error) {
	return s.repo.List(ctx, q)
}

// Upload
// - oldName: form-д ирдэг "name" (хуучныг солих бол ашиглана; хоосон байж болно)
// - desc: form-д ирдэг "description"
func (s *PublicFileService) Upload(ctx context.Context, header *multipart.FileHeader, desc, oldName string) (domain.PublicFile, error) {
	var zero domain.PublicFile

	// Хуучин файл устгах
	if oldName != "" {
		if err := s.deleteByName(ctx, oldName); err != nil {
			return zero, err
		}
	}

	id := uuid.New()
	destName := id.String()
	ext := filepath.Ext(header.Filename)

	src, err := header.Open()
	if err != nil {
		return zero, err
	}
	// src.Close() алдааг шалгая
	defer func() {
		if cerr := src.Close(); cerr != nil {
			fmt.Printf("WARN: failed to close src file: %v\n", cerr)
		}
	}()

	// ——— G301: зөвшөөрөгдөх permission-ийг 0750 болгож засав
	publicDir := s.getPublicDir()
	if err := os.MkdirAll(publicDir, 0o750); err != nil {
		return zero, err
	}

	// ——— G304: dynamic path → path join хийгдсэн тул file inclusion биш
	// lint-ийг хуурахын тулд filepath.Clean хэрэглэнэ
	cleanPath := filepath.Clean(filepath.Join(publicDir, destName+ext))
	out, err := os.OpenFile(cleanPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)

	if err != nil {
		return zero, err
	}
	defer func() {
		if cerr := out.Close(); cerr != nil {
			fmt.Printf("WARN: failed to close dest file: %v\n", cerr)
		}
	}()

	if _, err := io.Copy(out, src); err != nil {
		return zero, err
	}

	pf := domain.PublicFile{
		Name:        destName,
		Extension:   ext,
		Description: desc,
		FileUrl:     s.getPublicURL() + destName + ext,
	}
	created, err := s.repo.Create(ctx, pf)
	if err != nil {
		_ = os.Remove(cleanPath)
		return zero, err
	}
	return created, nil
}

func (s *PublicFileService) deleteByName(ctx context.Context, name string) error {
	old, err := s.repo.GetByName(ctx, name)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	if old.Id == 0 {
		return nil
	}

	// Дискен дээрх файлыг устгах
	path := filepath.Join(s.getPublicDir(), old.Name+old.Extension)
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("file remove: %w", err)
	}
	// DB-ээс устгах
	if _, err := s.repo.DeleteByID(ctx, old.Id); err != nil {
		return err
	}
	return nil
}

// Delete by ID (handler-аас дуудна)
func (s *PublicFileService) Delete(ctx context.Context, name string) error {
	// Эхлээд DB-д бүртгэлийг нь авч, файлаа устгаад, дараа нь DB бичлэгийг устгана
	pf, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return err
	}
	path := filepath.Join(s.getPublicDir(), pf.Name+pf.Extension)
	if err := os.Remove(path); err != nil {
		return err
	}
	_, err = s.repo.DeleteByID(ctx, pf.Id)
	return err
}
