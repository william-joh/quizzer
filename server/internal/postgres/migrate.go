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
	signup_date TIMESTAMP NOT NULL DEFAULT NOW(),
	calories_daily_goal INT DEFAULT 0,
	CONSTRAINT fk_language FOREIGN KEY(language) REFERENCES languages(language_name) ON DELETE CASCADE
);

CREATE TABLE sessions (
	id TEXT PRIMARY KEY,
	user_id TEXT,
	expiry_time INT NOT NULL,
	CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);
	`,
		`-- noop`)

	if err := m.Migrate(ctx); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	return nil
}
