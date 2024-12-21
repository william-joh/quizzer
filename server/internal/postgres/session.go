package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Session interface {
	Close() error
}

type db struct {
	pool *pgxpool.Pool
}

func (db *db) Close() error {
	db.pool.Close()
	return nil
}
