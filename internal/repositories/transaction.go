package repositories

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"infotecs-test-task/internal/models"
)

type TransactionRepository struct {
	Repository
}

func NewTransactionRepository(db *sqlx.DB) *TransactionRepository {
	return &TransactionRepository{Repository{Db: db}}
}

func (w *TransactionRepository) InsertOne(transaction models.Transaction) error {
	const fn = packagePath + "transaction.InsertOne"

	switch w.tx {
	case nil:
		query := "INSERT INTO transaction (\"to\", \"from\", amount, date) VALUES (:to, :from, :amount, :date)"
		_, err := w.Db.NamedExec(query, transaction)
		if err != nil {
			return fmt.Errorf("%s: exec statement: %w", fn, err)
		}
	default:
		query := "INSERT INTO transaction (\"to\", \"from\", amount, date) VALUES ($1, $2, $3, $4)"
		_, err := w.tx.Exec(query, transaction.To, transaction.From, transaction.Amount, transaction.Date)
		if err != nil {
			return fmt.Errorf("%s: exec statement: %w", fn, err)
		}
	}

	return nil
}

func (w *TransactionRepository) FindAllByWalletID(id string) ([]models.Transaction, error) {
	const fn = packagePath + "transaction.FindOneByID"

	transactions := make([]models.Transaction, 0)

	query := "SELECT date, \"to\", \"from\", amount FROM transaction WHERE \"to\"=$1 OR \"from\"=$1"

	switch w.tx {
	case nil:
		err := w.Db.Select(&transactions, query, id)
		if err != nil {
			return nil, fmt.Errorf("%s: exec query: %w", fn, err)
		}
	default:
		stmt, err := w.tx.Prepare(query)
		if err != nil {
			return nil, fmt.Errorf("%s: prepare statement: %w", fn, err)
		}

		rows, err := stmt.Query(id)
		if err != nil {
			return nil, fmt.Errorf("%s: exec statement: %w", fn, err)
		}

		for rows.Next() {
			var tr models.Transaction
			err = rows.Scan(&tr.Date, &tr.To, &tr.From, &tr.Amount)
			if err != nil {
				return nil, fmt.Errorf("%s: scan rows: %w", fn, err)
			}
			transactions = append(transactions, tr)
		}
	}

	return transactions, nil
}
