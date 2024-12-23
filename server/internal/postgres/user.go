package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/rs/zerolog/log"
	"github.com/william-joh/quizzer/server/internal/quizzer"
)

func (s *session) CreateUser(ctx context.Context, username, password string) error {
	log.Debug().Str("username", username).Msg("creating user")
	sql, args, err := psql().Insert("users").
		Columns("id", "username", "password").
		Values(sq.Expr("gen_random_uuid()"), username, password).ToSql()
	if err != nil {
		return err
	}

	_, err = s.conn.Exec(ctx, sql, args...)
	return err
}

func (s *session) GetUser(ctx context.Context, username string) (quizzer.User, error) {
	log.Debug().Str("username", username).Msg("getting user")

	sql, args, err := psql().Select("id", "username", "password", "signup_date").
		From("users").
		Where(sq.Eq{"username": username}).ToSql()
	if err != nil {
		return quizzer.User{}, err
	}

	row := s.conn.QueryRow(ctx, sql, args...)
	var user quizzer.User
	err = row.Scan(&user.ID, &user.Username, &user.Password, &user.SignupDate)
	return user, err
}

func (s *session) DeleteUser(ctx context.Context, username string) error {
	log.Debug().Str("username", username).Msg("deleting user")

	sql, args, err := psql().Delete("users").
		Where(sq.Eq{"username": username}).ToSql()
	if err != nil {
		return err
	}

	_, err = s.conn.Exec(ctx, sql, args...)
	return err
}
