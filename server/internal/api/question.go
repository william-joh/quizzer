package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/william-joh/quizzer/server/internal/quizzer"
)

func (s *server) createQuestionHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assert that the user is the creator of the quiz
		userID := r.Context().Value(userIDKey).(string)
		quizID := mux.Vars(r)["id"]
		quiz, err := s.db.Do(r.Context()).GetQuiz(r.Context(), quizID)
		if err != nil {
			log.Error().Err(err).Msg("failed to get quiz")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if quiz.CreatedBy != userID {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		var question quizzer.Question
		if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
			log.Error().Err(err).Msg("failed to decode request")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		question.ID = uuid.New().String()
		question.QuizID = mux.Vars(r)["id"]

		if err := s.db.Do(r.Context()).CreateQuestion(r.Context(), question); err != nil {
			log.Error().Err(err).Msg("failed to create question")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := struct {
			ID string `json:"id"`
		}{
			ID: question.ID,
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

func (s *server) listQuestionsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		quizID := mux.Vars(r)["id"]
		questions, err := s.db.Do(r.Context()).ListQuestions(r.Context(), quizID)
		if err != nil {
			log.Error().Err(err).Msg("failed to list questions")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(questions)
		if err != nil {
			log.Error().Err(err).Msg("failed to encode response")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func (s *server) deleteQuestionHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assert that the user is the creator of the quiz
		userID := r.Context().Value(userIDKey).(string)
		quizID := mux.Vars(r)["id"]
		quiz, err := s.db.Do(r.Context()).GetQuiz(r.Context(), quizID)
		if err != nil {
			log.Error().Err(err).Msg("failed to get quiz")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if quiz.CreatedBy != userID {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		questionID := mux.Vars(r)["questionID"]
		if err := s.db.Do(r.Context()).DeleteQuestion(r.Context(), questionID); err != nil {
			log.Error().Err(err).Msg("failed to delete question")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
