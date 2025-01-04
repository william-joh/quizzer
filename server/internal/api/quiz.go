package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/william-joh/quizzer/server/internal/postgres"
	"github.com/william-joh/quizzer/server/internal/quizzer"
)

func (s *server) createQuizHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var quiz struct {
			Title     string             `json:"title"`
			Questions []quizzer.Question `json:"questions"`
		}
		if err := json.NewDecoder(r.Body).Decode(&quiz); err != nil {
			log.Error().Err(err).Msg("failed to decode request body")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userID := r.Context().Value(userIDKey).(string)
		quizID := uuid.New().String()

		err := s.db.InTx(r.Context(), func(s postgres.Session) error {
			if err := s.CreateQuiz(r.Context(), quizID, quiz.Title, userID); err != nil {
				return err
			}

			for _, q := range quiz.Questions {
				q.ID = uuid.New().String()
				q.QuizID = quizID
				if err := s.CreateQuestion(r.Context(), q); err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to create quiz")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := struct {
			ID string `json:"id"`
		}{
			ID: quizID,
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Error().Err(err).Msg("failed to encode response")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func (s *server) getQuizHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		var quiz struct {
			Quiz      quizzer.Quiz       `json:"quiz"`
			Questions []quizzer.Question `json:"questions"`
		}

		err := s.db.InTx(r.Context(), func(s postgres.Session) error {
			q, err := s.GetQuiz(r.Context(), id)
			if err != nil {
				return err
			}
			quiz.Quiz = q

			questions, err := s.ListQuestions(r.Context(), id)
			if err != nil {
				return err
			}
			quiz.Questions = questions

			return nil
		})
		if err != nil {
			toJSONError(w, err, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(quiz)
		if err != nil {
			log.Error().Err(err).Msg("failed to encode response")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func (s *server) listQuizzesHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		quizzes, err := s.db.Do(r.Context()).ListQuizzes(r.Context())
		if err != nil {
			log.Error().Err(err).Msg("failed to list quizzes")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(quizzes)
		if err != nil {
			log.Error().Err(err).Msg("failed to encode response")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func (s *server) deleteQuizHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		// Get the quiz and check if the user is the owner
		quiz, err := s.db.Do(r.Context()).GetQuiz(r.Context(), id)
		if err != nil {
			log.Error().Err(err).Msg("failed to get quiz")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		userID := r.Context().Value(userIDKey).(string)
		if quiz.CreatedBy != userID {
			log.Error().Str("userID", userID).Str("createdBy", quiz.CreatedBy).Msg("unauthorized")
			http.Error(w, "access denied", http.StatusForbidden)
			return
		}

		if err := s.db.Do(r.Context()).DeleteQuiz(r.Context(), id); err != nil {
			log.Error().Err(err).Msg("failed to delete quiz")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
