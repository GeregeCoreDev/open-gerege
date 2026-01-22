// Package repository provides implementation for repository
//
// File: auth_repo.go
// Description: Repository for authentication-related entities
package repository

import (
	"context"
	"time"

	"templatev25/internal/domain"

	"gorm.io/gorm"
)

// AuthRepository defines the interface for authentication data access
type AuthRepository interface {
	// Credentials
	GetCredentialByUserID(ctx context.Context, userID int) (*domain.UserCredential, error)
	GetCredentialByEmail(ctx context.Context, email string) (*domain.UserCredential, error)
	CreateCredential(ctx context.Context, cred *domain.UserCredential) error
	UpdateCredential(ctx context.Context, cred *domain.UserCredential) error
	IncrementFailedAttempts(ctx context.Context, userID int) error
	ResetFailedAttempts(ctx context.Context, userID int) error
	LockAccount(ctx context.Context, userID int, until time.Time) error
	UnlockAccount(ctx context.Context, userID int) error

	// MFA TOTP
	GetMFAByUserID(ctx context.Context, userID int) (*domain.UserMFATotp, error)
	CreateMFA(ctx context.Context, mfa *domain.UserMFATotp) error
	UpdateMFA(ctx context.Context, mfa *domain.UserMFATotp) error
	DeleteMFA(ctx context.Context, userID int) error
	EnableMFA(ctx context.Context, userID int) error
	DisableMFA(ctx context.Context, userID int) error

	// MFA Backup Codes
	GetBackupCodes(ctx context.Context, userID int) ([]domain.UserMFABackupCode, error)
	GetUnusedBackupCodes(ctx context.Context, userID int) ([]domain.UserMFABackupCode, error)
	CreateBackupCodes(ctx context.Context, codes []domain.UserMFABackupCode) error
	DeleteBackupCodes(ctx context.Context, userID int) error
	UseBackupCode(ctx context.Context, codeID int) error

	// Sessions (DB layer - Redis is primary)
	CreateSession(ctx context.Context, session *domain.Session) error
	GetSession(ctx context.Context, id string) (*domain.Session, error)
	GetUserSessions(ctx context.Context, userID int) ([]domain.Session, error)
	GetActiveUserSessions(ctx context.Context, userID int) ([]domain.Session, error)
	UpdateSessionActivity(ctx context.Context, id string) error
	RevokeSession(ctx context.Context, id string, reason string) error
	RevokeAllUserSessions(ctx context.Context, userID int, reason string) error

	// Login History
	CreateLoginHistory(ctx context.Context, history *domain.LoginHistory) error
	GetLoginHistory(ctx context.Context, userID int, limit int) ([]domain.LoginHistory, error)
	GetRecentLoginHistory(ctx context.Context, userID int, since time.Time) ([]domain.LoginHistory, error)

	// Security Audit Trail
	CreateAuditTrail(ctx context.Context, audit *domain.SecurityAuditTrail) error
	GetAuditTrail(ctx context.Context, userID int, limit int) ([]domain.SecurityAuditTrail, error)
	GetAuditTrailByAction(ctx context.Context, userID int, action string, limit int) ([]domain.SecurityAuditTrail, error)

	// Password History
	GetPasswordHistory(ctx context.Context, userID int, limit int) ([]domain.PasswordHistory, error)
	CreatePasswordHistory(ctx context.Context, history *domain.PasswordHistory) error

	// User Status
	UpdateUserStatus(ctx context.Context, userID int, status string, reason string, changedBy int) error
	UpdateUserLoginStats(ctx context.Context, userID int) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type authRepository struct {
	db *gorm.DB
}

// NewAuthRepository creates a new auth repository instance
func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

// ============================================================
// CREDENTIALS
// ============================================================

func (r *authRepository) GetCredentialByUserID(ctx context.Context, userID int) (*domain.UserCredential, error) {
	var cred domain.UserCredential
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&cred).Error
	if err != nil {
		return nil, err
	}
	return &cred, nil
}

func (r *authRepository) GetCredentialByEmail(ctx context.Context, email string) (*domain.UserCredential, error) {
	var cred domain.UserCredential
	err := r.db.WithContext(ctx).
		Joins("JOIN users ON users.id = user_credentials.user_id").
		Where("users.email = ? AND users.deleted_date IS NULL", email).
		First(&cred).Error
	if err != nil {
		return nil, err
	}
	return &cred, nil
}

func (r *authRepository) CreateCredential(ctx context.Context, cred *domain.UserCredential) error {
	return r.db.WithContext(ctx).Create(cred).Error
}

func (r *authRepository) UpdateCredential(ctx context.Context, cred *domain.UserCredential) error {
	return r.db.WithContext(ctx).Save(cred).Error
}

func (r *authRepository) IncrementFailedAttempts(ctx context.Context, userID int) error {
	return r.db.WithContext(ctx).
		Model(&domain.UserCredential{}).
		Where("user_id = ?", userID).
		Update("failed_login_attempts", gorm.Expr("failed_login_attempts + 1")).Error
}

func (r *authRepository) ResetFailedAttempts(ctx context.Context, userID int) error {
	return r.db.WithContext(ctx).
		Model(&domain.UserCredential{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"failed_login_attempts": 0,
			"locked_until":          nil,
		}).Error
}

func (r *authRepository) LockAccount(ctx context.Context, userID int, until time.Time) error {
	return r.db.WithContext(ctx).
		Model(&domain.UserCredential{}).
		Where("user_id = ?", userID).
		Update("locked_until", until).Error
}

func (r *authRepository) UnlockAccount(ctx context.Context, userID int) error {
	return r.db.WithContext(ctx).
		Model(&domain.UserCredential{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"failed_login_attempts": 0,
			"locked_until":          nil,
		}).Error
}

// ============================================================
// MFA TOTP
// ============================================================

func (r *authRepository) GetMFAByUserID(ctx context.Context, userID int) (*domain.UserMFATotp, error) {
	var mfa domain.UserMFATotp
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&mfa).Error
	if err != nil {
		return nil, err
	}
	return &mfa, nil
}

func (r *authRepository) CreateMFA(ctx context.Context, mfa *domain.UserMFATotp) error {
	return r.db.WithContext(ctx).Create(mfa).Error
}

func (r *authRepository) UpdateMFA(ctx context.Context, mfa *domain.UserMFATotp) error {
	return r.db.WithContext(ctx).Save(mfa).Error
}

func (r *authRepository) DeleteMFA(ctx context.Context, userID int) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&domain.UserMFATotp{}).Error
}

func (r *authRepository) EnableMFA(ctx context.Context, userID int) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&domain.UserMFATotp{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"is_enabled":  true,
			"verified_at": now,
		}).Error
}

func (r *authRepository) DisableMFA(ctx context.Context, userID int) error {
	return r.db.WithContext(ctx).
		Model(&domain.UserMFATotp{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"is_enabled":  false,
			"verified_at": nil,
		}).Error
}

// ============================================================
// MFA BACKUP CODES
// ============================================================

func (r *authRepository) GetBackupCodes(ctx context.Context, userID int) ([]domain.UserMFABackupCode, error) {
	var codes []domain.UserMFABackupCode
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_date DESC").
		Find(&codes).Error
	return codes, err
}

func (r *authRepository) GetUnusedBackupCodes(ctx context.Context, userID int) ([]domain.UserMFABackupCode, error) {
	var codes []domain.UserMFABackupCode
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND used_at IS NULL", userID).
		Find(&codes).Error
	return codes, err
}

func (r *authRepository) CreateBackupCodes(ctx context.Context, codes []domain.UserMFABackupCode) error {
	return r.db.WithContext(ctx).Create(&codes).Error
}

func (r *authRepository) DeleteBackupCodes(ctx context.Context, userID int) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&domain.UserMFABackupCode{}).Error
}

func (r *authRepository) UseBackupCode(ctx context.Context, codeID int) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&domain.UserMFABackupCode{}).
		Where("id = ?", codeID).
		Update("used_at", now).Error
}

// ============================================================
// SESSIONS
// ============================================================

func (r *authRepository) CreateSession(ctx context.Context, session *domain.Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *authRepository) GetSession(ctx context.Context, id string) (*domain.Session, error) {
	var session domain.Session
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *authRepository) GetUserSessions(ctx context.Context, userID int) ([]domain.Session, error) {
	var sessions []domain.Session
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_date DESC").
		Find(&sessions).Error
	return sessions, err
}

func (r *authRepository) GetActiveUserSessions(ctx context.Context, userID int) ([]domain.Session, error) {
	var sessions []domain.Session
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND revoked_at IS NULL AND expires_at > ?", userID, time.Now()).
		Order("created_date DESC").
		Find(&sessions).Error
	return sessions, err
}

func (r *authRepository) UpdateSessionActivity(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&domain.Session{}).
		Where("id = ?", id).
		Update("last_activity_at", time.Now()).Error
}

func (r *authRepository) RevokeSession(ctx context.Context, id string, reason string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&domain.Session{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"revoked_at":     now,
			"revoked_reason": reason,
		}).Error
}

func (r *authRepository) RevokeAllUserSessions(ctx context.Context, userID int, reason string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&domain.Session{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Updates(map[string]interface{}{
			"revoked_at":     now,
			"revoked_reason": reason,
		}).Error
}

// ============================================================
// LOGIN HISTORY
// ============================================================

func (r *authRepository) CreateLoginHistory(ctx context.Context, history *domain.LoginHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *authRepository) GetLoginHistory(ctx context.Context, userID int, limit int) ([]domain.LoginHistory, error) {
	var history []domain.LoginHistory
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_date DESC").
		Limit(limit).
		Find(&history).Error
	return history, err
}

func (r *authRepository) GetRecentLoginHistory(ctx context.Context, userID int, since time.Time) ([]domain.LoginHistory, error) {
	var history []domain.LoginHistory
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND created_date >= ?", userID, since).
		Order("created_date DESC").
		Find(&history).Error
	return history, err
}

// ============================================================
// SECURITY AUDIT TRAIL
// ============================================================

func (r *authRepository) CreateAuditTrail(ctx context.Context, audit *domain.SecurityAuditTrail) error {
	return r.db.WithContext(ctx).Create(audit).Error
}

func (r *authRepository) GetAuditTrail(ctx context.Context, userID int, limit int) ([]domain.SecurityAuditTrail, error) {
	var audit []domain.SecurityAuditTrail
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_date DESC").
		Limit(limit).
		Find(&audit).Error
	return audit, err
}

func (r *authRepository) GetAuditTrailByAction(ctx context.Context, userID int, action string, limit int) ([]domain.SecurityAuditTrail, error) {
	var audit []domain.SecurityAuditTrail
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND action = ?", userID, action).
		Order("created_date DESC").
		Limit(limit).
		Find(&audit).Error
	return audit, err
}

// ============================================================
// PASSWORD HISTORY
// ============================================================

func (r *authRepository) GetPasswordHistory(ctx context.Context, userID int, limit int) ([]domain.PasswordHistory, error) {
	var history []domain.PasswordHistory
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_date DESC").
		Limit(limit).
		Find(&history).Error
	return history, err
}

func (r *authRepository) CreatePasswordHistory(ctx context.Context, history *domain.PasswordHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

// ============================================================
// USER STATUS
// ============================================================

func (r *authRepository) UpdateUserStatus(ctx context.Context, userID int, status string, reason string, changedBy int) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"status":            status,
			"status_reason":     reason,
			"status_changed_at": now,
			"status_changed_by": changedBy,
		}).Error
}

func (r *authRepository) UpdateUserLoginStats(ctx context.Context, userID int) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"last_login_at": now,
			"login_count":   gorm.Expr("login_count + 1"),
		}).Error
}

func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Where("email = ? AND deleted_date IS NULL", email).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
