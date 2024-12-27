package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type contextKey string

const userIDKey contextKey = "userID"

func (s *server) authHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			log.Error().Err(err).Msg("failed to decode request body")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := s.db.Do(r.Context()).GetUserByUsername(r.Context(), creds.Username)
		if err != nil {
			log.Error().Err(err).Msg("failed to get user")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
			log.Error().Err(err).Msg("failed to compare password")
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		sessionID := uuid.New().String()
		s.db.Do(r.Context()).CreateAuthSession(r.Context(), user.ID, sessionID)

		// set cookie
		http.SetCookie(w, &http.Cookie{
			Name:  "quizzer_session_id",
			Value: sessionID,
			// HttpOnly: true,
			// SameSite: http.SameSiteStrictMode,
			// Secure: true,
		})
		w.WriteHeader(http.StatusOK)
	})
}

func (s *server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("quizzer_session_id")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := s.db.Do(r.Context()).GetAuthSession(r.Context(), sessionID.Value)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
