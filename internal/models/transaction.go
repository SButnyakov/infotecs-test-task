package models

import (
	"time"
)

type Transaction struct {
	Id     uint64    `db:"id"`
	To     string    `db:"to"`
	From   string    `db:"from"`
	Amount float32   `db:"amount"`
	Date   time.Time `db:"date"`
}
