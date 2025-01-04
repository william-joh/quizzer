package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
)

const versionTable = "schema_version_default"

func migrateDb(ctx context.Context, conn *pgx.Conn) error {
	m, err := migrate.NewMigrator(ctx, conn, versionTable)
	if err != nil {
		return fmt.Errorf("new migrator: %w", err)
	}

	m.AppendMigration("initial setup",
		`
CREATE TABLE users (
	id TEXT PRIMARY KEY,
	username TEXT UNIQUE NOT NULL,
	password TEXT NOT NULL,
	signup_date TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE auth_sessions (
	id TEXT PRIMARY KEY,
	user_id TEXT,
	expiry_time INT NOT NULL,
	CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE quizzes (
	id TEXT PRIMARY KEY,
	title TEXT NOT NULL,
	created_by TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	CONSTRAINT fk_user FOREIGN KEY(created_by) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE questions (
	id TEXT PRIMARY KEY,
	quiz_id TEXT NOT NULL,
	question TEXT NOT NULL,
	index INT NOT NULL,
	time_limit_seconds INT NOT NULL,
	answers TEXT[] NOT NULL CHECK (array_length(answers, 1) > 1),
	correct_answers TEXT[] NOT NULL CHECK (array_length(correct_answers, 1) > 0 AND correct_answers <@ answers),
	video_url TEXT,
	video_start_time_seconds INT,
	video_end_time_seconds INT,
	CONSTRAINT fk_quiz FOREIGN KEY(quiz_id) REFERENCES quizzes(id) ON DELETE CASCADE,
	CONSTRAINT unique_quiz_index UNIQUE (quiz_id, index)
);
	`,
		`-- noop`)

	if err := m.Migrate(ctx); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	return nil
}
