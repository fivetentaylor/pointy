package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/fivetentaylor/pointy/pkg/config"
	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/query"
)

type Manager struct {
	jwt    *JWTEngine
	google *GoogleOauth
	email  *Email
	query  *query.Query
	env    string
	secret []byte
}

func NewManager(cfg config.Server) *Manager {
	return &Manager{
		jwt:    NewJWT(cfg.JWTSecret),
		google: NewGoogle(cfg.GoogleOauth),
		email:  NewEmail(),
		env:    os.Getenv("ENV"),
		secret: []byte(cfg.JWTSecret),
	}
}

func (m *Manager) Routes(r chi.Router) {
	r.Use(m.JWTMiddleware)
	r.Get("/auth/signout", m.Signout)
	r.Get("/auth/google", m.GoogleRedirect)
	r.Get("/auth/google/callback", m.GoogleCallback)
	r.Post("/auth/login", m.email.EmailLogin)
	r.Post("/auth/signup", m.email.EmailSignup)
	r.Post("/auth/token", OneTimeAccessLink)
	r.Post("/auth/magic_link", SendMagicLink)
	r.Post("/auth/refresh_token", RefreshToken)

	r.Get("/login", m.GetLogin)
	r.Get("/access/{code}", GetMagicLink)
	r.Get("/invite/{code}", m.GetInvite)
}

func (m *Manager) AttachUserClaimIfExists(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cookie, err := r.Cookie(constants.CookieName)
		if err == nil && cookie != nil {
			claims, err := m.jwt.ParseUserToken(cookie.Value)
			if err == nil {
				ctx := env.UserClaimCtx(r.Context(), claims)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			log.Error(fmt.Errorf("error parsing jwt: %w", err))
		}

		// Via Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := m.jwt.ParseUserToken(token)
			if err == nil {
				ctx = env.UserClaimCtx(r.Context(), claims)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			log.Error(fmt.Errorf("error parsing jwt: %w", err))
		}

		log.Info("no user claim found")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Manager) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := env.UserClaim(r.Context())
		if err != nil {
			log.Error(fmt.Errorf("error getting admin user: %w", err))
			http.NotFound(w, r)
			return
		}

		if !user.Admin {
			log.Error(fmt.Errorf("user is not admin"))
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Manager) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := m.jwt.Attach(r.Context())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Manager) JWT() *JWTEngine {
	return m.jwt
}

func (m *Manager) Signout(w http.ResponseWriter, r *http.Request) {
	expiredCookie := &http.Cookie{
		Name:     constants.CookieName,
		Domain:   os.Getenv("COOKIE_DOMAIN"),
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	// Set the new cookie, which will effectively remove the old one
	http.SetCookie(w, expiredCookie)

	render.JSON(w, r, map[string]string{"status": "ok"})
}
