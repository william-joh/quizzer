package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/william-joh/quizzer/server/internal/quizzer"
)

type Session interface {
	CreateUser(ctx context.Context, username, password string) error
	GetUser(ctx context.Context, username string) (quizzer.User, error)
	DeleteUser(ctx context.Context, username string) error
	Close() error
}

type db struct {
	pool *pgxpool.Pool
}

func (db *db) Close() error {
	log.Debug().Msg("closing db...")
	db.pool.Close()
	log.Debug().Msg("db closed")
	return nil
}

func psql() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}
