package main

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/william-joh/quizzer/server/internal/api"
	"github.com/william-joh/quizzer/server/internal/execution"
	"github.com/william-joh/quizzer/server/internal/postgres"
)

func main() {
	ctx := context.Background()

	db, err := postgres.Connect(ctx)
	if err != nil {
		log.Panic().Err(err).Msg("failed to setup db")
	}
	defer db.Close()

	executioner := execution.NewInMemory(db)

	api := api.NewAPI(db, executioner)
	api.Run()
}
