package worker

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jpoz/conveyor"
	cwire "github.com/jpoz/conveyor/wire"
	"github.com/posthog/posthog-go"
	"github.com/redis/go-redis/v9"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"

	"github.com/teamreviso/code/pkg/background/jobs"
	"github.com/teamreviso/code/pkg/client"
	"github.com/teamreviso/code/pkg/config"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/graph/loaders"
	"github.com/teamreviso/code/pkg/pubsub"
	"github.com/teamreviso/code/pkg/query"
	"github.com/teamreviso/code/pkg/storage/dynamo"
	"github.com/teamreviso/code/pkg/storage/s3"
	"github.com/teamreviso/code/pkg/utils"
)

type Worker struct {
	Conveyor *conveyor.Worker

	cancel        context.CancelFunc
	cfg           config.Worker
	ctx           context.Context
	dynamo        *dynamo.DB
	log           *log.Logger
	slog          *slog.Logger
	openai        *openai.Client
	pubsub        pubsub.PubSubInterface
	query         *query.Query
	redis         *redis.Client
	s3            *s3.S3
	ses           *client.SES
	freeplay      *client.Freeplay
	rawdb         *gorm.DB
	bg            *conveyor.Client
	posthogClient posthog.Client
	loaders       *loaders.Loaders
}

func New(cfg config.Worker, db *gorm.DB, r *redis.Client) (*Worker, context.CancelFunc, error) {
	workerCtx, cancel := context.WithCancel(context.Background())

	cfg.Logger = utils.NewSlogFromEnv()

	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
	}).WithPrefix("server")

	level, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err == nil {
		logger.SetLevel(level)
	}

	openAiClient := openai.NewClient(cfg.OpenAIKey)

	dydb, err := dynamo.NewDB()
	if err != nil {
		cancel()
		return nil, nil, err
	}

	bg, err := conveyor.NewClient(&cfg)
	if err != nil {
		cancel()
		return nil, nil, err
	}

	s3, err := s3.NewS3()
	if err != nil {
		cancel()
		return nil, nil, err
	}

	ses, err := client.NewSESFromEnv(bg)
	if err != nil {
		cancel()
		return nil, nil, err
	}

	q := query.Use(db)

	worker, err := conveyor.NewWorker(&cfg)
	if err != nil {
		cancel()
		return nil, nil, err
	}

	phClient, err := posthog.NewWithConfig(
		os.Getenv("PUBLIC_POSTHOG_KEY"),
		posthog.Config{
			PersonalApiKey: os.Getenv("POSTHOG_SERVER_FEATURE_FLAG_KEY"), // Optional, but much more performant.  If this token is not supplied, then fetching feature flag values will be slower.
			Endpoint:       os.Getenv("PUBLIC_POSTHOG_HOST"),
		},
	)
	if err != nil {
		cancel()
		return nil, nil, err
	}

	// Initialize loaders
	l := loaders.NewLoaders()

	w := &Worker{
		Conveyor: worker,

		cancel:        cancel,
		cfg:           cfg,
		ctx:           workerCtx,
		log:           logger,
		slog:          cfg.Logger,
		openai:        openAiClient,
		ses:           ses,
		dynamo:        dydb,
		query:         q,
		redis:         r,
		s3:            s3,
		rawdb:         db,
		bg:            bg,
		posthogClient: phClient,
		loaders:       l,
	}

	w.setup()

	return w, cancel, nil
}

func (w *Worker) setup() {
	w.Conveyor.Use(SentryMiddleware)
	w.Conveyor.Use(w.Attach)
	w.Conveyor.Use(func(ctx context.Context, job *cwire.Job, next func(ctx context.Context) error) error {
		log := env.SLog(ctx)

		log.Info("starting job", "jid", job.Uuid, "queue", job.Queue, "type", job.Type)
		err := next(ctx)
		log.Info("finished job", "jid", job.Uuid, "queue", job.Queue, "type", job.Type, "error", err)

		return err
	})
	w.Conveyor.RegisterJobs(jobs.AllJobs...)
}

func (w *Worker) Attach(ctx context.Context, job *cwire.Job, next func(ctx context.Context) error) error {
	ctx = env.Attach(
		ctx,
		w.log,
		w.slog,
		w.query,
		w.openai,
		w.dynamo,
		w.ses,
		w.redis,
		w.posthogClient,
		w.bg,
		w.s3,
		w.rawdb,
	)
	// Add loaders to the context using the existing loadersKey
	ctx = context.WithValue(ctx, loaders.LoadersKey, w.loaders)
	return next(ctx)
}

func (w *Worker) Run() error {
	w.log.Info("running worker", "addr", w.cfg.Addr)
	err := w.Conveyor.Run(w.ctx)
	w.log.Info("worker stopped", "error", err)
	return err
}

func (w *Worker) RunHealthServer() error {
	r := chi.NewRouter()
	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health check
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{\"status\": \"ok\"}"))
	})

	w.log.Info("running worker health server", "addr", w.cfg.Addr)
	return http.ListenAndServe(w.cfg.Addr, r)
}

func (w *Worker) Close() {
	w.log.Warn("Close called on worker")
	w.cancel()
	<-w.ctx.Done()
}
