package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/rs/zerolog/log"
	"github.com/william-joh/quizzer/server/internal/quizzer"
)

func (s *session) CreateQuiz(ctx context.Context, id, title, createdBy string) error {
	log.Debug().Str("id", id).Str("title", title).Str("createdBy", createdBy).Msg("creating quiz")
	sql, args, err := psql().Insert("quizzes").
		Columns("id", "title", "created_by").
		Values(id, title, createdBy).ToSql()
	if err != nil {
		return err
	}

	_, err = s.conn.Exec(ctx, sql, args...)
	return err
}

func (s *session) GetQuiz(ctx context.Context, id string) (quizzer.Quiz, error) {
	log.Debug().Str("id", id).Msg("getting quiz")

	sql, args, err := psql().Select("id", "title", "created_by", "created_at").
		From("quizzes").
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return quizzer.Quiz{}, err
	}

	row := s.conn.QueryRow(ctx, sql, args...)
	var quiz quizzer.Quiz
	err = row.Scan(&quiz.ID, &quiz.Title, &quiz.CreatedBy, &quiz.CreatedAt)
	return quiz, err
}

func (s *session) DeleteQuiz(ctx context.Context, id string) error {
	log.Debug().Str("id", id).Msg("deleting quiz")

	sql, args, err := psql().Delete("quizzes").
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}

	_, err = s.conn.Exec(ctx, sql, args...)
	return err
}

func (s *session) ListQuizzes(ctx context.Context) ([]quizzer.Quiz, error) {
	log.Debug().Msg("listing quizzes")

	sql, args, err := psql().Select("id", "title", "created_by", "created_at").
		From("quizzes").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := s.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quizzes []quizzer.Quiz
	for rows.Next() {
		var quiz quizzer.Quiz
		if err := rows.Scan(&quiz.ID, &quiz.Title, &quiz.CreatedBy, &quiz.CreatedAt); err != nil {
			return nil, err
		}
		quizzes = append(quizzes, quiz)
	}

	return quizzes, nil
}
