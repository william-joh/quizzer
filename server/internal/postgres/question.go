package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/rs/zerolog/log"
	"github.com/william-joh/quizzer/server/internal/quizzer"
)

func (s *session) CreateQuestion(ctx context.Context, question quizzer.Question) error {
	log.Debug().Str("id", question.ID).Str("quizID", question.QuizID).Str("question", question.Question).Int("index", question.Index).Dur("timeLimit", question.TimeLimit).Strs("answers", question.Answers).Str("correctAnswer", question.CorrectAnswer).Msg("creating question")

	sql, args, err := psql().Insert("questions").
		Columns("id", "quiz_id", "question", "index", "time_limit", "answers", "correct_answer", "video_url", "video_start_time", "video_end_time").
		Values(question.ID, question.QuizID, question.Question, question.Index, question.TimeLimit, question.Answers, question.CorrectAnswer, question.VideoURL, question.VideoStartTime, question.VideoEndTime).ToSql()
	if err != nil {
		return err
	}

	_, err = s.conn.Exec(ctx, sql, args...)
	return err
}

func (s *session) GetQuestion(ctx context.Context, id string) (quizzer.Question, error) {
	log.Debug().Str("id", id).Msg("getting question")

	sql, args, err := psql().Select("id", "quiz_id", "question", "index", "time_limit", "answers", "correct_answer", "video_url", "video_start_time", "video_end_time").
		From("questions").
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return quizzer.Question{}, err
	}

	row := s.conn.QueryRow(ctx, sql, args...)
	var question quizzer.Question
	err = row.Scan(&question.ID, &question.QuizID, &question.Question, &question.Index, &question.TimeLimit, &question.Answers, &question.CorrectAnswer, &question.VideoURL, &question.VideoStartTime, &question.VideoEndTime)
	return question, err
}

func (s *session) ListQuestions(ctx context.Context, quizID string) ([]quizzer.Question, error) {
	log.Debug().Str("quizID", quizID).Msg("listing questions")

	sql, args, err := psql().Select("id", "quiz_id", "question", "index", "time_limit", "answers", "correct_answer", "video_url", "video_start_time", "video_end_time").
		From("questions").
		Where(sq.Eq{"quiz_id": quizID}).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := s.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []quizzer.Question
	for rows.Next() {
		var question quizzer.Question
		err = rows.Scan(&question.ID, &question.QuizID, &question.Question, &question.Index, &question.TimeLimit, &question.Answers, &question.CorrectAnswer, &question.VideoURL, &question.VideoStartTime, &question.VideoEndTime)
		if err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}
	return questions, nil
}

func (s *session) UpdateQuestion(ctx context.Context, question quizzer.Question) error {
	log.Debug().Str("id", question.ID).Str("question", question.Question).Int("index", question.Index).Dur("timeLimit", question.TimeLimit).Strs("answers", question.Answers).Str("correctAnswer", question.CorrectAnswer).Msg("updating question")

	sql, args, err := psql().Update("questions").
		SetMap(map[string]interface{}{
			"question":         question.Question,
			"index":            question.Index,
			"time_limit":       question.TimeLimit,
			"answers":          question.Answers,
			"correct_answer":   question.CorrectAnswer,
			"video_url":        question.VideoURL,
			"video_start_time": question.VideoStartTime,
			"video_end_time":   question.VideoEndTime,
		}).
		Where(sq.Eq{"id": question.ID}).ToSql()
	if err != nil {
		return err
	}

	_, err = s.conn.Exec(ctx, sql, args...)
	return err
}

func (s *session) DeleteQuestion(ctx context.Context, id string) error {
	log.Debug().Str("id", id).Msg("deleting question")

	sql, args, err := psql().Delete("questions").
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}

	_, err = s.conn.Exec(ctx, sql, args...)
	return err
}
