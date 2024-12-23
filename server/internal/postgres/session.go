package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/william-joh/quizzer/server/internal/quizzer"
)

type Conn interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Database interface {
	Do(ctx context.Context) Session
	InTx(ctx context.Context, fn func(Session) error) error
	Close() error
}

type db struct {
	pool *pgxpool.Pool
}

func (db *db) Do(ctx context.Context) Session {
	return &session{conn: db.pool}
}

func (db *db) InTx(ctx context.Context, fn func(Session) error) error {
	tx, err := db.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	if err := fn(&session{conn: tx}); err != nil {
		defer tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (db *db) Close() error {
	log.Debug().Msg("closing db...")
	db.pool.Close()
	log.Debug().Msg("db closed")
	return nil
}

type Session interface {
	CreateUser(ctx context.Context, username, password string) error
	GetUser(ctx context.Context, username string) (quizzer.User, error)
	DeleteUser(ctx context.Context, username string) error
}

var _ Session = &session{}

type session struct {
	conn Conn
}

func psql() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}
