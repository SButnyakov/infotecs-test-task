package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

const packagePath = "internal.repositories."

type TX interface {
	GetTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	SetTx(tx *sql.Tx)
	Commit() error
	Rollback() error
}

type Repository struct {
	Db *sqlx.DB
	tx *sql.Tx
}

func (r *Repository) GetTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	if r.tx != nil {
		return r.tx, nil
	}
	tx, err := r.Db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	r.tx = tx
	return r.tx, nil
}

func (r *Repository) SetTx(tx *sql.Tx) {
	r.tx = tx
}

func (r *Repository) Commit() error {
	if r.tx == nil {
		return fmt.Errorf("unable to commit nil transaction")
	}
	return r.tx.Commit()
}

func (r *Repository) Rollback() error {
	if r.tx == nil {
		return fmt.Errorf("unable to rollback nil transaction")
	}
	return r.tx.Commit()
}
