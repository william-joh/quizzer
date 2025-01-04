package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func (s *server) createUserHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		userID := uuid.New().String()
		err = s.db.Do(r.Context()).CreateUser(r.Context(), userID, user.Username, string(hashedPassword))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sessionID := uuid.New().String()
		s.db.Do(r.Context()).CreateAuthSession(r.Context(), userID, sessionID)

		// set cookie
		http.SetCookie(w, &http.Cookie{
			Name:  "quizzer_session_id",
			Value: sessionID,
			// HttpOnly: true,
			// SameSite: http.SameSiteStrictMode,
			// Secure: true,
		})

		resp := struct {
			ID string `json:"id"`
		}{
			ID: userID,
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func (s *server) getUserHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(userIDKey).(string)
		user, err := s.db.Do(r.Context()).GetUser(r.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := struct {
			ID         string `json:"id"`
			Username   string `json:"username"`
			SignupDate string `json:"signupDate"`
		}{
			ID:         user.ID,
			Username:   user.Username,
			SignupDate: user.SignupDate.Format("2006-01-02T15:04:05Z"),
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func (s *server) deleteUserHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		err := s.db.Do(r.Context()).DeleteUser(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
