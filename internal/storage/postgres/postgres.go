package postgres

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"infotecs-test-task/internal/lib/config"
)

func NewSQL(cfg config.PG) (*sql.DB, error) {
	url := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name)

	conn, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func NewSQLX(cfg config.PG) (*sqlx.DB, error) {
	db, err := NewSQL(cfg)
	if err != nil {
		return nil, err
	}
	return sqlx.NewDb(db, "postgres"), nil
}
