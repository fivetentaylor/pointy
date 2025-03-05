package config

import (
	"log/slog"

	"github.com/teamreviso/code/pkg/constants"
)

type Server struct {
	Addr            string
	ImageTag        string
	GoogleOauth     *GoogleOauth
	Logger          *slog.Logger
	JWTSecret       string
	OpenAIKey       string
	AllowedOrigins  []string
	WorkerRedisAddr string
}

type GoogleOauth struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// implement github.com/jpoz/conveyor/config.ClientConfig
func (s Server) GetLogger() *slog.Logger {
	if s.Logger == nil {
		return slog.Default()
	}
	return s.Logger
}

func (s *Server) SetLogger(logger *slog.Logger) {
	s.Logger = logger
}

func (s Server) GetRedisURL() string {
	return s.WorkerRedisAddr
}

func (s Server) GetNamespace() string {
	return constants.DefaultJobNamespace
}
