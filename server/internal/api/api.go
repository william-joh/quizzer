package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/william-joh/quizzer/server/internal/postgres"
)

type API interface {
	Handler() http.Handler
}

type server struct {
	db postgres.Database
}

func NewAPI(db postgres.Database) API {
	return &server{db: db}
}

func (s *server) Handler() http.Handler {
	r := mux.NewRouter()
	r.Use(loggerMiddleware)

	r.Handle("/auth", s.authHandler()).Methods(http.MethodPost)
	r.Handle("/users", s.createUserHandler()).Methods(http.MethodPost)
	r.HandleFunc("/healthz", healthzHandler)

	authorized := r.NewRoute().Subrouter()
	authorized.Use(s.authMiddleware)

	authorized.Handle("/users/{id}", s.getUserHandler()).Methods(http.MethodGet)
	authorized.Handle("/users/{id}", s.deleteUserHandler()).Methods(http.MethodDelete)

	authorized.Handle("/quizzes", s.createQuizHandler()).Methods(http.MethodPost)
	authorized.Handle("/quizzes", s.listQuizzesHandler()).Methods(http.MethodGet)
	authorized.Handle("/quizzes/{id}", s.getQuizHandler()).Methods(http.MethodGet)
	authorized.Handle("/quizzes/{id}", s.deleteQuizHandler()).Methods(http.MethodDelete)

	authorized.Handle("/quizzes/{id}/questions", s.createQuestionHandler()).Methods(http.MethodPost)
	authorized.Handle("/quizzes/{id}/questions", s.listQuestionsHandler()).Methods(http.MethodGet)
	authorized.Handle("/quizzes/{id}/questions/{questionID}", s.deleteQuestionHandler()).Methods(http.MethodDelete)

	return r
}

// func (s *server) Run() error {
// 	srv := &http.Server{
// 		Handler: r,
// 		Addr:    "127.0.0.1:8000",
// 		// Good practice: enforce timeouts for servers you create!
// 		WriteTimeout: 15 * time.Second,
// 		ReadTimeout:  15 * time.Second,
// 	}

// 	return srv.ListenAndServe()
// }

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Str("method", r.Method).Str("url", r.URL.String()).Msg("request")
		next.ServeHTTP(w, r)
	})
}
