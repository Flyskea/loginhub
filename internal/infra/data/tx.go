package data

import (
	"context"

	"gorm.io/gorm"

	"loginhub/internal/base/iface"
)

var _ iface.Transaction = (*TXManager)(nil)

type contextTxKey struct{}

type TXManager struct {
	db *gorm.DB
}

func NewTXManager(db *gorm.DB) *TXManager {
	return &TXManager{
		db: db,
	}
}

func (t *TXManager) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return t.db.WithContext(ctx)
}

func (t *TXManager) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if !ok {
		tx = t.db
	}
	return tx.Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}
