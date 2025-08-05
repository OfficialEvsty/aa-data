package db

import (
	"context"
	"database/sql"
	"fmt"
)

type txKeyType string

const txKey txKeyType = "tx"

// TxManager provides logic with tx working
type TxManager struct {
	exec *sql.DB
}

func WithTxInContext(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func TxFromContext(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txKey).(*sql.Tx)
	return tx, ok
}

func NewTxManager(db *sql.DB) *TxManager {
	tx := &TxManager{exec: db}
	if db == nil {
		panic("NewServerImporter: db is nil because tx is nil") // 💥 остановись сразу
	}
	return tx
}

// WithTx working with tx life-circle
func (m *TxManager) WithTx(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
	tx, ok := TxFromContext(ctx)
	if ok {
		err := fn(ctx, tx)
		if err != nil {
			return err
		}
		return nil
	}
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
