package models

type Wallet struct {
	Id      string  `db:"id"`
	Balance float32 `db:"balance"`
}
