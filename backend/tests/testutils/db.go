// Package testutils provides test utilities for integration tests
//
// File: db.go
// Description: Test database helpers
package testutils

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	testDB   *gorm.DB
	dbOnce   sync.Once
	dbMutex  sync.Mutex
	tablesMu sync.Mutex
)

// DefaultTestDSN returns the default test database DSN
func DefaultTestDSN() string {
	if dsn := os.Getenv("TEST_DB_DSN"); dsn != "" {
		return dsn
	}
	return "postgres://test:test@localhost:5433/test_db?sslmode=disable"
}

// SetupTestDB creates or returns a shared test database connection.
// The connection is reused across tests for performance.
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dbOnce.Do(func() {
		dsn := DefaultTestDSN()
		var err error

		// Retry connection with backoff
		for i := 0; i < 5; i++ {
			testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Silent),
			})
			if err == nil {
				break
			}
			time.Sleep(time.Second * time.Duration(i+1))
		}

		if err != nil {
			panic(fmt.Sprintf("failed to connect to test database: %v", err))
		}

		// Configure connection pool
		sqlDB, err := testDB.DB()
		if err == nil {
			sqlDB.SetMaxIdleConns(5)
			sqlDB.SetMaxOpenConns(10)
			sqlDB.SetConnMaxLifetime(time.Hour)
		}
	})

	if testDB == nil {
		t.Fatal("test database not initialized")
	}

	return testDB
}

// SetupTestDBWithTx returns a database wrapped in a transaction.
// The transaction is rolled back at the end of the test for isolation.
func SetupTestDBWithTx(t *testing.T) *gorm.DB {
	t.Helper()

	db := SetupTestDB(t)

	// Start a transaction
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}

	// Rollback at the end of the test
	t.Cleanup(func() {
		tx.Rollback()
	})

	return tx
}

// TruncateTable truncates the specified table.
// Use with caution - this deletes all data.
func TruncateTable(t *testing.T, db *gorm.DB, tableName string) {
	t.Helper()

	tablesMu.Lock()
	defer tablesMu.Unlock()

	if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", tableName)).Error; err != nil {
		t.Logf("warning: failed to truncate table %s: %v", tableName, err)
	}
}

// TruncateTables truncates multiple tables.
func TruncateTables(t *testing.T, db *gorm.DB, tableNames ...string) {
	t.Helper()

	for _, table := range tableNames {
		TruncateTable(t, db, table)
	}
}

// ResetSequence resets the ID sequence for a table.
func ResetSequence(t *testing.T, db *gorm.DB, tableName string) {
	t.Helper()

	query := fmt.Sprintf("ALTER SEQUENCE %s_id_seq RESTART WITH 1", tableName)
	if err := db.Exec(query).Error; err != nil {
		t.Logf("warning: failed to reset sequence for %s: %v", tableName, err)
	}
}

// CleanupTestDB closes the test database connection.
// Call this in TestMain if needed.
func CleanupTestDB() {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	if testDB != nil {
		if sqlDB, err := testDB.DB(); err == nil {
			sqlDB.Close()
		}
		testDB = nil
	}
}
