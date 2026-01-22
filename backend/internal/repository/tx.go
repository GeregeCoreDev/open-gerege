// Package repository provides implementation for repository
//
// File: tx.go
// Description: implementation for repository
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package repository

import (
	"context"

	"gorm.io/gorm"
)

type TxFunc func(tx *gorm.DB) error

// WithTx runs fn within a transaction and propagates the given ctx.
// It guarantees ctx cancellation/timeout will abort DB ops.
func WithTx(ctx context.Context, db *gorm.DB, fn TxFunc) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx.WithContext(ctx))
	})
}
