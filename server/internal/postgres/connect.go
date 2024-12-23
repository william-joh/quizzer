package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func Connect(ctx context.Context) (Session, error) {
	// postgres://username:password@localhost:5432
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("must specify DATABASE_URL environment variable")
	}

	log.Debug().Str("db_url", dbURL).Msg("connecting to db")
	dbpool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("connect to db: %w", err)
	}

	log.Debug().Msg("connected to db")
	conn, err := dbpool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	if err := migrateDb(ctx, conn.Conn()); err != nil {
		return nil, fmt.Errorf("migrate db: %w", err)
	}

	return &db{pool: dbpool}, nil
}
