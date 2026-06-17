package sqlite

import (
	"context"
	"database/sql"

	"github.com/raiyin/artserver/internal/port"
)

// Tx wraps sql.Tx to implement port.Tx.
type Tx struct {
	*sql.Tx
}

// Commit commits the transaction.
func (t *Tx) Commit() error {
	return t.Tx.Commit()
}

// Rollback rolls back the transaction.
func (t *Tx) Rollback() error {
	return t.Tx.Rollback()
}

// TransactionManager implements port.TransactionManager.
type TransactionManager struct {
	db *sql.DB
}

// NewTransactionManager creates a new TransactionManager.
func NewTransactionManager(db *sql.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

// BeginTx begins a new transaction.
func (tm *TransactionManager) BeginTx(ctx context.Context) (port.Tx, error) {
	tx, err := tm.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Tx{Tx: tx}, nil
}

// Ensure interfaces are satisfied.
var _ port.TransactionManager = (*TransactionManager)(nil)
var _ port.Tx = (*Tx)(nil)
