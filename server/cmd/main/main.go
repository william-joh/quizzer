package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/william-joh/quizzer/server/internal/postgres"
)

func main() {
	ctx := context.Background()

	db, err := postgres.Connect(ctx)
	if err != nil {
		log.Panic().Err(err).Msg("failed to setup db")
	}
	defer db.Close()

	startServer()
}

func startServer() {
	log.Debug().Msg("Starting server...")

	r := mux.NewRouter()

	r.HandleFunc("/healthz", HealthzHandler)

	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
}

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}
