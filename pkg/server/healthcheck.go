package server

import (
	"net/http"

	"github.com/teamreviso/code/pkg/env"
)

func (s *Server) Healthcheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := env.Log(ctx)

		log.Info("Healthcheck", "status", "ok")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok", "version": "` + s.Config.ImageTag + `"}`))
	}
}
