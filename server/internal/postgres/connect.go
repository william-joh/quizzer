package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context) (Session, error) {
	// postgres://username:password@localhost:5432/database_name
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("must specify DATABASE_URL environment variable")
	}

	dbpool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("connect to db: %w", err)
	}

	conn, err := dbpool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Conn().Close(ctx)

	if err := migrateDb(ctx, conn.Conn()); err != nil {
		return nil, fmt.Errorf("migrate db: %w", err)
	}

	return &db{pool: dbpool}, nil
}
