package api

import (
	"context"
	"encoding/json"
	"fmt"
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
			toJSONError(w, fmt.Errorf("decode request body: %w", err), http.StatusBadRequest)
			return
		}

		user, err := s.db.Do(r.Context()).GetUserByUsername(r.Context(), creds.Username)
		if err != nil {
			toJSONError(w, fmt.Errorf("get user: %w", err), http.StatusBadRequest)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
			toJSONError(w, fmt.Errorf("check password: %w", err), http.StatusUnauthorized)
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
			toJSONError(w, fmt.Errorf("get cookie: %w", err), http.StatusUnauthorized)
			return
		}

		userID, err := s.db.Do(r.Context()).GetAuthSession(r.Context(), sessionID.Value)
		if err != nil {
			toJSONError(w, fmt.Errorf("get auth session: %w", err), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type errorResponse struct {
	Error string `json:"error"`
}

func toJSONError(w http.ResponseWriter, err error, status int) {
	log.Error().Err(err).Msg("responding with error")

	errorResponse := errorResponse{Error: err.Error()}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		log.Error().Err(err).Msg("failed to encode error response")
	}
}
