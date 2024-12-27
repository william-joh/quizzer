package postgres

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/rs/zerolog/log"
)

const sessionValidSeconds = 60 * 60 * 24 // 24 days

func (s *session) CreateAuthSession(ctx context.Context, userID, sessionID string) error {
	log.Debug().Str("userID", userID).Msg("creating auth session")

	sql, args, err := psql().Insert("auth_sessions").
		Columns("id", "user_id", "expiry_time").
		Values(sessionID, userID, time.Now().Unix()+int64(sessionValidSeconds)).ToSql()
	if err != nil {
		return err
	}

	_, err = s.conn.Exec(ctx, sql, args...)
	return err
}

func (s *session) GetAuthSession(ctx context.Context, sessionID string) (string, error) {
	log.Debug().Str("sessionID", sessionID).Msg("getting auth session")

	sql, args, err := psql().Select("user_id", "expiry_time").
		From("auth_sessions").
		Where(sq.Eq{"id": sessionID}).ToSql()
	if err != nil {
		return "", err
	}

	row := s.conn.QueryRow(ctx, sql, args...)
	var userID string
	var expiryTime int64
	err = row.Scan(&userID, &expiryTime)
	if err != nil {
		return "", err
	}

	return userID, err
}
