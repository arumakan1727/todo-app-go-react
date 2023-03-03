package pgsql

import (
	"context"
	"fmt"

	"github.com/arumakan1727/todo-app-go-react/domain"
)

func (r *repository) BeginTx(ctx context.Context) error {
	if r.tx_internal != nil {
		return nil
	}
	tx, err := r.db_internal.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	r.tx_internal = tx
	r.db = tx
	return nil
}

func (r *repository) CommitTx(context.Context) error {
	tx := r.tx_internal
	if tx == nil {
		return fmt.Errorf("cannot commit: %w", domain.ErrNotInTransaction)
	}
	r.tx_internal = nil
	r.db = r.db_internal

	err := tx.Commit()
	return err
}

func (r *repository) RollbackTx(context.Context) error {
	tx := r.tx_internal
	if tx == nil {
		return fmt.Errorf("cannot rollback: %w", domain.ErrNotInTransaction)
	}
	r.tx_internal = nil
	r.db = r.db_internal

	err := tx.Rollback()
	return err
}
