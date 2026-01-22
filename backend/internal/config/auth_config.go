// Package config provides local configuration for auth and related features
//
// File: auth_config.go
// Description: Configuration for local authentication, Redis, and security settings
package config

import (
	"os"
	"strconv"
	"time"
)

// RedisConfig holds Redis connection settings
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// Addr returns the Redis address in host:port format
func (c *RedisConfig) Addr() string {
	return c.Host + ":" + c.Port
}

// LocalAuthConfig holds local authentication settings
type LocalAuthConfig struct {
	// Enabled indicates if local authentication is enabled
	Enabled bool

	// SessionTTL is the session lifetime
	SessionTTL time.Duration

	// MFATokenTTL is the MFA pending token lifetime
	MFATokenTTL time.Duration

	// LockoutThreshold is the number of failed attempts before lockout
	LockoutThreshold int

	// LockoutDuration is how long the account stays locked
	LockoutDuration time.Duration

	// PasswordMinLength is the minimum password length
	PasswordMinLength int

	// PasswordHistoryCount is how many previous passwords to check
	PasswordHistoryCount int

	// TOTPIssuer is the issuer name for TOTP QR codes
	TOTPIssuer string

	// EncryptionKey is the 32-byte key for encrypting TOTP secrets
	EncryptionKey string
}

// AuthConfig combines all auth-related configurations
type AuthConfig struct {
	Redis     RedisConfig
	LocalAuth LocalAuthConfig
}

// LoadAuthConfig loads authentication configuration from environment variables
func LoadAuthConfig() *AuthConfig {
	return &AuthConfig{
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		LocalAuth: LocalAuthConfig{
			Enabled:              getEnvBool("LOCAL_AUTH_ENABLED", true),
			SessionTTL:           getEnvDuration("LOCAL_AUTH_SESSION_TTL", 24*time.Hour),
			MFATokenTTL:          getEnvDuration("LOCAL_AUTH_MFA_TOKEN_TTL", 5*time.Minute),
			LockoutThreshold:     getEnvInt("LOCAL_AUTH_LOCKOUT_THRESHOLD", 5),
			LockoutDuration:      getEnvDuration("LOCAL_AUTH_LOCKOUT_DURATION", 15*time.Minute),
			PasswordMinLength:    getEnvInt("LOCAL_AUTH_PASSWORD_MIN_LENGTH", 8),
			PasswordHistoryCount: getEnvInt("LOCAL_AUTH_PASSWORD_HISTORY_COUNT", 5),
			TOTPIssuer:           getEnv("LOCAL_AUTH_TOTP_ISSUER", "TemplateBackend"),
			EncryptionKey:        getEnv("LOCAL_AUTH_ENCRYPTION_KEY", ""),
		},
	}
}

// getEnv returns the environment variable value or a default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt returns the environment variable as int or a default
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// getEnvBool returns the environment variable as bool or a default
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

// getEnvDuration returns the environment variable as duration or a default
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
