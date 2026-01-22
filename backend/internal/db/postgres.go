// Package db provides implementation for db
//
// File: postgres.go
// Description: implementation for db
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
/*
Package db нь database connection-ийг удирдана.

PostgreSQL database-тай GORM ORM ашиглан холбогдоно.

Features:
  - Connection pooling (MaxIdleConns, MaxOpenConns)
  - Connection lifetime management
  - Schema prefix (table naming)
  - SQL logging

Ашиглалт:

	gormDB, err := db.NewPostgres(cfg)
	if err != nil {
	    log.Fatal("db connection failed", zap.Error(err))
	}
	defer func() {
	    sqlDB, _ := gormDB.DB()
	    sqlDB.Close()
	}()
*/
package db

import (
	"fmt"  // String formatting
	"time" // Duration

	"git.gerege.mn/backend-packages/config" // Configuration

	"gorm.io/driver/postgres" // PostgreSQL driver
	"gorm.io/gorm"            // ORM
	"gorm.io/gorm/logger"     // SQL logging
	"gorm.io/gorm/schema"     // Table naming
)

// ============================================================
// NEW POSTGRES
// ============================================================

// NewPostgres нь PostgreSQL database-д GORM connection үүсгэнэ.
//
// Parameters:
//   - cfg: Application configuration (DB host, port, user, password, etc.)
//
// Returns:
//   - *gorm.DB: GORM database instance
//   - error: Connection error
//
// Connection string format:
//
//	host=localhost port=5432 user=postgres password=secret dbname=template sslmode=disable
//
// GORM Config:
//   - Logger: SQL query-г лог хийнэ
//   - NamingStrategy: Table prefix (schema.table_name)
//
// Connection Pool:
//   - MaxIdleConns: Idle connection-уудын max тоо
//   - MaxOpenConns: Open connection-уудын max тоо
//   - ConnMaxLifetime: Connection-ийн max lifetime
//   - ConnMaxIdleTime: Idle connection-ийн max хугацаа
//
// Жишээ:
//
//	gormDB, err := db.NewPostgres(cfg)
//	if err != nil {
//	    log.Fatal("db init failed", zap.Error(err))
//	}
//
//	// Application shutdown хийхэд
//	sqlDB, _ := gormDB.DB()
//	sqlDB.Close()
func NewPostgres(cfg config.Config) (*gorm.DB, error) {
	// ============================================================
	// STEP 1: DSN (Data Source Name) үүсгэх
	// ============================================================
	// PostgreSQL connection string format
	// search_path нэмэх: schema name-ийг connection-д тохируулна
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=%s",
		cfg.DB.Host,     // Database server host
		cfg.DB.Port,     // Database server port (default: 5432)
		cfg.DB.User,     // Database user
		cfg.DB.Password, // Database password
		cfg.DB.Name,     // Database name
		cfg.DB.Schema,   // Schema name (search_path)
	)

	// ============================================================
	// STEP 2: GORM connection үүсгэх
	// ============================================================
	g, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// SQL logging (query-г console-д харуулна)
		// Production-д logger.Silent ашиглаж болно
		Logger: logger.Default.LogMode(logger.Info),

		// Table naming strategy
		// Schema prefix: "template_backend.users" гэх мэт
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: fmt.Sprintf("%s.", cfg.DB.Schema),
		},
	})
	if err != nil {
		return nil, err
	}

	// ============================================================
	// STEP 3: Connection pool тохируулах
	// ============================================================
	// *sql.DB instance авах (underlying database/sql connection)
	sqlDB, err := g.DB()
	if err != nil {
		return nil, err
	}

	// MaxIdleConns: Idle state-д байж болох max connection тоо
	// Хэтэрхий их бол санах ой идэнэ
	// Хэтэрхий бага бол шинэ connection үүсгэх хугацаа алдана
	// Зөвлөмж: MaxOpenConns-ийн 10-20%
	sqlDB.SetMaxIdleConns(cfg.DB.MaxIdleConn)

	// MaxOpenConns: Нийт open connection-ийн max тоо
	// Database server-ийн max_connections-оос бага байх ёстой
	// Concurrent request тоонд тохируулна
	sqlDB.SetMaxOpenConns(cfg.DB.MaxOpenConn)

	// ConnMaxLifetime: Connection-ийн max амьдрах хугацаа
	// Энэ хугацааны дараа connection хаагдаж, шинээр үүснэ
	// Database server-ийн timeout-оос бага байх ёстой
	sqlDB.SetConnMaxLifetime(cfg.DB.MaxConnLifetime)

	// ConnMaxIdleTime: Idle connection-ийн max хугацаа
	// Idle connection энэ хугацааны дараа хаагдана
	// Resource cleanup-д тусална
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)

	return g, nil
}
