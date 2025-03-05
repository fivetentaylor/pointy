package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/charmbracelet/log"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
	"github.com/jpoz/conveyor"
	"github.com/posthog/posthog-go"
	"github.com/redis/go-redis/v9"
	"github.com/riandyrn/otelchi"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"

	"github.com/teamreviso/code/pkg/admin"
	"github.com/teamreviso/code/pkg/assets"
	"github.com/teamreviso/code/pkg/client"
	"github.com/teamreviso/code/pkg/config"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/graph"
	"github.com/teamreviso/code/pkg/query"
	"github.com/teamreviso/code/pkg/server/auth"
	"github.com/teamreviso/code/pkg/server/waitlist"
	"github.com/teamreviso/code/pkg/storage/dynamo"
	sRedis "github.com/teamreviso/code/pkg/storage/redis"
	"github.com/teamreviso/code/pkg/storage/s3"
	"github.com/teamreviso/code/pkg/utils"
	"github.com/teamreviso/code/pkg/views"
)

func NewServer(
	cfg config.Server,
	db *gorm.DB,
) (*Server, error) {
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
	}).WithPrefix("server")

	slogger := utils.NewSlogFromEnv()

	level, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err == nil {
		logger.SetLevel(level)
	}

	openAiClient := openai.NewClient(cfg.OpenAIKey)

	dynamo, err := dynamo.NewDB()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to dynamodb: %w", err)
	}

	bg, err := conveyor.NewClient(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to background worker: %w", err)
	}

	s3, err := s3.NewS3()
	if err != nil {
		return nil, err
	}

	// Redis
	rc, err := sRedis.NewRedis()
	if err != nil {
		return nil, err
	}

	ses, err := client.NewSESFromEnv(bg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to email: %w", err)
	}

	q := query.Use(db)

	sentryMiddleware := sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	})

	phClient, err := posthog.NewWithConfig(
		os.Getenv("PUBLIC_POSTHOG_KEY"),
		posthog.Config{
			PersonalApiKey: os.Getenv("POSTHOG_SERVER_FEATURE_FLAG_KEY"), // Optional, but much more performant.  If this token is not supplied, then fetching feature flag values will be slower.
			Endpoint:       os.Getenv("PUBLIC_POSTHOG_HOST"),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to posthog: %w", err)
	}
	defer phClient.Close()

	return &Server{
		Config:           cfg,
		Auth:             auth.NewManager(cfg),
		Query:            q,
		RawDB:            db,
		SES:              ses,
		OpenAi:           openAiClient,
		Dynamo:           dynamo,
		Redis:            rc,
		Posthog:          phClient,
		Background:       bg,
		Logger:           logger,
		SLog:             slogger,
		S3:               s3,
		SentryMiddleware: sentryMiddleware,
	}, nil
}

type Server struct {
	Config config.Server
	Auth   *auth.Manager

	Logger           *log.Logger
	SLog             *slog.Logger
	Query            *query.Query
	SES              *client.SES
	OpenAi           *openai.Client
	Dynamo           *dynamo.DB
	Redis            *redis.Client
	Posthog          posthog.Client
	Background       *conveyor.Client
	S3               *s3.S3
	RawDB            *gorm.DB
	SentryMiddleware *sentryhttp.Handler
	HttpServer       *http.Server
	Listener         net.Listener
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // For the sake of example, allow all origins
	},
}

func (s *Server) Router() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(Recoverer) // custom recoverer using slog
	r.Use(s.SentryMiddleware.Handle)

	// Attach
	r.Use(s.EnvMiddleware)

	// User
	r.Use(s.Auth.AttachUserClaimIfExists)

	// OpenTelemetry tracing to honeycomb
	r.Use(otelchi.Middleware("reviso-server", otelchi.WithChiRoutes(r)))

	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   s.Config.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Cookie", "sentry-trace", "baggage"},
		ExposedHeaders:   []string{"Link", "SetCookie"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
		Debug:            (log.GetLevel() == log.DebugLevel),
	}))

	// / will redrirect to /login
	r.Get("/", func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/login", http.StatusFound) })

	// Health check
	r.Get("/healthcheck", s.Healthcheck())

	// Static
	r.Get("/static/*", assets.Static("pkg/assets/static", "/static/", os.Getenv("ENV") == "production"))
	r.Get("/src/*", assets.SrcHandler("/src"))

	// UI
	s.UI(r)

	// Auth
	r.Group(s.Auth.Routes)
	r.Route("/waitlist", func(r chi.Router) {
		r.Get("/success", waitlist.WaitlistSuccess)
		r.Post("/", waitlist.AddToWaitlist)
	})

	// Routes
	r.Route("/graphql", func(r chi.Router) {
		r.Get("/", playground.Handler("GraphQL playground", "/graphql/query"))
		r.Mount("/query", graph.NewHandler())
	})

	// Api
	r.Route("/api/v1", func(r chi.Router) {
		r.HandleFunc("/documents/{docID}/rogue/ws", s.RogueWebSocket)
		r.HandleFunc("/documents/{docID}/threads/{threadID}/authors/{authorID}/stream", s.StreamingVoice)
		r.Get("/documents/{docID}/doc.html", s.HtmlDocument)
		r.Get("/documents/{docID}/editor.html", s.DocumentEditor)
		r.Post("/documents/{docID}/images/", s.CreateDocumentImage)
		r.Get("/documents/{docID}/images/{imageID}", s.GetDocumentImage)
		r.Get("/users/{userID}/avatar", s.GetUserAvatar)
		r.Put("/avatar", s.UpdateUserAvatar)
	})

	// Admin
	r.Route("/admin", func(r chi.Router) {
		r.Use(s.Auth.RequireAdmin)
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := env.RawDBCtx(r.Context(), s.RawDB)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		})

		r.Group(admin.Routes)
	})

	// Checkout
	r.Route("/payments", func(r chi.Router) {
		r.Post("/checkout", s.Checkout)
		r.Get("/cancel", s.CancelCheckout)
		r.Get("/success", s.CheckoutSuccessful)
		r.Get("/status", s.CheckoutStatus)
		r.Get("/failure", s.CheckoutFailure)
	})

	r.Post("/stripe/webhook", s.StripeWebhook)

	// Not found
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		views.NotFound().Render(context.Background(), w)
	})

	if os.Getenv("ENV") == "development" {
		r.HandleFunc("/dev/preview_email/{email}", s.PreviewEmail)
		r.HandleFunc("/dev/send_email/{email}", s.EmailMe)
		r.HandleFunc("/dev/err", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("test")
		}))
	}

	if os.Getenv("ENV") == "test" || true {
		r.Route("/test", func(r chi.Router) {
			r.Use(s.Auth.JWTMiddleware)
			r.Post("/documents", s.CreateTestDocument)
			r.Get("/documents/{docID}.txt", s.TestDocumentText)
			r.Get("/documents/{docID}.html", s.TestDocumentHTML)
			r.Get("/documents/{docID}", s.ViewTestDocument)
		})
	}

	return r
}

func (s *Server) ListenAndServe() error {
	if s.HttpServer != nil {
		return fmt.Errorf("server already running!")
	}
	s.Logger.Infof("Server running on %s", s.Config.Addr)
	router := s.Router()

	// Write all routes in debug mode
	chi.Walk(router, func(method string, route string, _ http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		s.Logger.Debugf("[%s]:\t%s has %d middlewares", method, route, len(middlewares))
		return nil
	})

	// This is to help with testability
	// We break out the HttpServer and Listener to allow for testing to grab the addr of the server listening on :0
	var err error
	s.HttpServer = &http.Server{Addr: s.Config.Addr, Handler: router}
	s.Listener, err = net.Listen("tcp", s.Config.Addr)
	if err != nil {
		return err
	}

	if os.Getenv("ENV") == "development" {
		/*sslCertDirArg := flag.String("d", "./dev/certs", "the directory of SSL cert")
		sslCrtNameArg := flag.String("c", "_wildcard.reviso.dev.pem", "the filename of SSL cert")
		sslKeyNameArg := flag.String("k", "_wildcard.reviso.dev-key.pem", "the filename of SSL key")
		flag.Parse()

		if string((*sslCrtNameArg)[0]) != "/" {
			*sslCrtNameArg = "/" + *sslCrtNameArg
		}
		if string((*sslKeyNameArg)[0]) != "/" {
			*sslKeyNameArg = "/" + *sslKeyNameArg
		}
		revisoCrtPath := *sslCertDirArg + *sslCrtNameArg
		sslCertKeyPath := *sslCertDirArg + *sslKeyNameArg
		*/
		revisoCrtPath := "./dev/certs/_wildcard.reviso.dev.pem"
		revisoKeyPath := "./dev/certs/_wildcard.reviso.dev-key.pem"
		revisoCert, err := tls.LoadX509KeyPair(revisoCrtPath, revisoKeyPath)
		if err != nil {
			return err
		}

		pointyCrtPath := "./dev/certs/_wildcard.dev.pointy.ai.pem"
		pointyKeyPath := "./dev/certs/_wildcard.dev.pointy.ai-key.pem"
		pointyCert, err := tls.LoadX509KeyPair(pointyCrtPath, pointyKeyPath)
		if err != nil {
			return err
		}

		tlsConfig := &tls.Config{
			GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
				fmt.Printf("TAYZEE ServerName: %q\n", info.ServerName)
				// Logic to select the right certificate based on the requested hostname
				if strings.HasSuffix(info.ServerName, ".reviso.dev") {
					return &revisoCert, nil
				} else if strings.HasSuffix(info.ServerName, ".dev.pointy.ai") {
					return &pointyCert, nil
				}
				// Default certificate
				return &revisoCert, nil
			},
		}

		s.HttpServer.TLSConfig = tlsConfig

		s.Logger.Warn("Server running in development https mode")
		return s.HttpServer.ServeTLS(s.Listener, "", "")
	}

	return s.HttpServer.Serve(s.Listener)
}
