package postgres_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/william-joh/quizzer/server/internal/postgres"
)

func SetupTestDB(t *testing.T) postgres.Database {
	connString := "postgres://postgres:mysecretpassword@localhost:5432"

	// create database
	dbName := "testdb"
	dbpool, err := pgxpool.New(context.Background(), connString)
	require.NoError(t, err)

	// drop test database if exists
	_, err = dbpool.Exec(context.Background(), fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
	require.NoError(t, err)

	// create test database
	_, err = dbpool.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", dbName))
	require.NoError(t, err)
	dbpool.Close()

	os.Setenv("DATABASE_URL", connString+"/"+dbName)
	db, err := postgres.Connect(context.Background())
	require.NoError(t, err)

	t.Cleanup(func() {
		fmt.Println("closing connection")
		db.Close()
	})

	return db
}
