package db

import (
	"context"
	"database/sql"
	"fmt"
)

// TxManager provides logic with tx working
type TxManager struct {
	exec *sql.DB
}

func NewTxManager(db *sql.DB) *TxManager {
	tx := &TxManager{exec: db}
	if tx == nil {
		panic("NewServerImporter: tx is nil") // 💥 остановись сразу
	}
	return tx
}

// WithTx working with tx life-circle
func (m *TxManager) WithTx(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
	tx, err := m.exec.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error while tx begin: %v", err)
	}
	err = fn(ctx, tx)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("error while tx rollback: %v", err)
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error while tx commit: %v", err)
	}
	return nil
}
