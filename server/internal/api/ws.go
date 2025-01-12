package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

func (s *server) wsHandler() http.Handler {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins in development
		},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := mux.Vars(r)["code"]
		log.Info().Str("code", code).Msg("WebSocket connection")

		e, err := s.exectioner.GetExecution(r.Context(), code)
		if err != nil {
			toJSONError(w, err, http.StatusInternalServerError)
			return
		}

		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error().Err(err).Msg("failed to upgrade connection to WebSocket")
			return
		}
		defer conn.Close()

		// Handle the WebSocket connection
		for {
			err := e.HandleMessages(conn)
			if err != nil {
				log.Error().Err(err).Msg("failed to handle message")
				// TODO: Handle error
				break
			}
		}
	})
}
