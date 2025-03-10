package server

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/fivetentaylor/pointy/pkg/env"
)

func (s *Server) EnvMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := middleware.GetReqID(r.Context())
		logger := s.Logger.With("req_id", reqID)
		slogger := s.SLog.With("req_id", reqID)
		ctx := env.Attach(r.Context(), logger, slogger, s.Query, s.OpenAi, s.Dynamo, s.SES, s.Redis, s.Posthog, s.Background, s.S3, s.RawDB)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
