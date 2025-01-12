package api

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
	"github.com/william-joh/quizzer/server/internal/execution"
	"github.com/william-joh/quizzer/server/internal/postgres"
)

type API interface {
	Handler() http.Handler
	Run() error
}

type server struct {
	db         postgres.Database
	exectioner execution.Service
}

func NewAPI(db postgres.Database, executioner execution.Service) API {
	return &server{db: db, exectioner: executioner}
}

func (s *server) Handler() http.Handler {
	r := mux.NewRouter()
	r.Use(loggerMiddleware)

	r.Handle("/game/{code}", s.wsHandler()).Methods(http.MethodGet)

	r.Handle("/auth", s.authHandler()).Methods(http.MethodPost)
	r.Handle("/users", s.createUserHandler()).Methods(http.MethodPost)
	r.HandleFunc("/healthz", healthzHandler)

	authorized := r.NewRoute().Subrouter()
	authorized.Use(s.authMiddleware)

	authorized.Handle("/current-user", s.getUserHandler()).Methods(http.MethodGet)
	authorized.Handle("/users/{id}", s.deleteUserHandler()).Methods(http.MethodDelete)

	authorized.Handle("/quizzes", s.createQuizHandler()).Methods(http.MethodPost)
	authorized.Handle("/quizzes/{id}/start", s.startQuizHandler()).Methods(http.MethodPost)
	authorized.Handle("/quizzes", s.listQuizzesHandler()).Methods(http.MethodGet)
	authorized.Handle("/quizzes/{id}", s.getQuizHandler()).Methods(http.MethodGet)
	authorized.Handle("/quizzes/{id}", s.deleteQuizHandler()).Methods(http.MethodDelete)

	return r
}

func (s *server) Run() error {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173"}, // All origins
		AllowedMethods:   []string{"GET", "POST", "DELETE", "HEAD"},
		AllowCredentials: true,
	})

	srv := &http.Server{
		Handler: c.Handler(s.Handler()),
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Info().Str("addr", srv.Addr).Msg("starting server")
	return srv.ListenAndServe()
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Str("method", r.Method).Str("url", r.URL.String()).Msg("request")
		next.ServeHTTP(w, r)
	})
}
