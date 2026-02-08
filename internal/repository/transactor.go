package repository

import (
	"context"

	"gorm.io/gorm"
)

type gormTransactor struct {
	db *gorm.DB
}

func NewGormTransactor(db *gorm.DB) *gormTransactor {
	return &gormTransactor{db: db}
}

// WithinTransaction ejecuta una función dentro de una transacción de GORM.
// Inyecta la transacción en el contexto para que los repositorios la encuentren.
func (t *gormTransactor) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		ctxWithTx := context.WithValue(ctx, "tx", tx)
		return fn(ctxWithTx)
	})
}
