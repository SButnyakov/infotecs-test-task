package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"infotecs-test-task/internal/models"
	"infotecs-test-task/internal/storage"
)

type WalletRepository struct {
	Repository
}

func NewWalletRepository(db *sqlx.DB) *WalletRepository {
	return &WalletRepository{Repository{Db: db}}
}

func (w *WalletRepository) InsertOne(wallet models.Wallet) error {
	const fn = packagePath + "wallet.InsertOne"

	var err error

	switch w.tx {
	case nil:
		_, err = w.Db.NamedExec("INSERT INTO wallet (id, balance) VALUES (:id, :balance)", wallet)
	default:
		_, err = w.tx.Exec("INSERT INTO wallet (id, balance) VALUES ($1, $2)", wallet.Id, wallet.Balance)
	}

	if err != nil {
		return fmt.Errorf("%s: exec query: %w", fn, err)
	}

	return nil
}

func (w *WalletRepository) UpdateBalance(id string, delta float32) error {
	const fn = packagePath + "wallet.UpdateBalance"

	query := "UPDATE wallet SET balance=balance+$1 WHERE id=$2"

	var res sql.Result
	var err error

	switch w.tx {
	case nil:
		res, err = w.Db.Exec(query, delta, id)
	default:
		res, err = w.tx.Exec(query, delta, id)
	}

	if err != nil {
		return fmt.Errorf("%s: exec statement: %w", fn, err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return storage.ErrWalletNotFound
	}

	return nil
}

func (w *WalletRepository) FindOneByID(id string) (models.Wallet, error) {
	const fn = packagePath + "wallet.FindOneByID"

	var wallet models.Wallet
	var err error

	switch w.tx {
	case nil:
		err = w.Db.Get(&wallet, "SELECT id, balance FROM wallet WHERE id=$1", id)
		if errors.Is(err, sql.ErrNoRows) {
			return wallet, storage.ErrWalletNotFound
		}
		if err != nil {
			return models.Wallet{}, fmt.Errorf("%s: exec statement: %w", fn, err)
		}
	default:
		row := w.tx.QueryRow("SELECT id, balance FROM wallet WHERE id=$1", id)
		if row == nil {
			return models.Wallet{}, storage.ErrWalletNotFound
		}
		err = row.Scan(&wallet.Id, &wallet.Balance)
		if err != nil {
			return models.Wallet{}, fmt.Errorf("%s: scan row: %w", fn, err)
		}
	}

	return wallet, nil
}
