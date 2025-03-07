package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/getsentry/sentry-go"
	"github.com/honeycombio/otel-config-go/otelconfig"
	"github.com/jpoz/conveyor"
	convconfig "github.com/jpoz/conveyor/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/fivetentaylor/pointy/pkg/assets"
	"github.com/fivetentaylor/pointy/pkg/background/wire"
	"github.com/fivetentaylor/pointy/pkg/background/worker"
	"github.com/fivetentaylor/pointy/pkg/config"
	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/server"
	"github.com/fivetentaylor/pointy/pkg/utils"
)

// This will be set by -ldflags in the Dockerfile to the current git tag
var ImageTag = "development"

func main() {
	var (
		serverFlag    = flag.Bool("server", false, "Run the server only (Also can set REVISO_MODE=server)")
		workerFlag    = flag.Bool("worker", false, "Run the worker only (Also can set REVISO_MODE=worker)")
		buildOnlyFlag = flag.Bool("build", false, "Only build the frontend")
		helpFlag      = flag.Bool("h", false, "Show help message")
	)
	flag.Parse()

	// Check if help was requested
	if *helpFlag {
		printHelp()
		os.Exit(0)
	}

	log := utils.NewSlogFromEnv()
	utils.SetDefaultLogger()

	revisoMode := os.Getenv("REVISO_MODE")
	runBoth := !*serverFlag && !*workerFlag && revisoMode == ""
	runServer := (*serverFlag || revisoMode == "server") || runBoth
	runWorker := (*workerFlag || revisoMode == "worker") || runBoth

	if *buildOnlyFlag {
		fmt.Printf("[reviso] Running build only\n")

		fmt.Printf(
			"[reviso] Building frontend: NODE_ENV=%s, API_HOST=%s, WS_HOST=%s\n, IMAGE_TAG=%s\n",
			os.Getenv("NODE_ENV"),
			os.Getenv("API_HOST"),
			os.Getenv("WS_HOST"),
			os.Getenv("IMAGE_TAG"),
		)
		err := assets.BuildAssets()
		if err != nil {
			log.Error("failed to build assets", "error", err)
			os.Exit(1)
		}

		return
	}

	fmt.Printf("[reviso] Running server: %t\n", runServer)
	fmt.Printf("[reviso] Running worker: %t\n", runWorker)

	gormdb := connectGorm()
	rc := connectToRedis(log)
	connectToSentry(log)

	otelShutdown, err := otelconfig.ConfigureOpenTelemetry()
	if err != nil {
		log.Error("error setting up OTel SDK", "error", err)
		os.Exit(1)
	}
	defer otelShutdown()
	defer sentry.Flush(2 * time.Second)
	defer rc.Close()

	completed := make(chan struct{})
	var workerCancel context.CancelFunc

	if runServer {
		allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
		if allowedOrigins == "" {
			allowedOrigins = "https://*,http://*"
		}

		cfg := config.Server{
			Addr:           os.Getenv("ADDR"),
			ImageTag:       ImageTag,
			JWTSecret:      os.Getenv("JWT_SECRET"),
			OpenAIKey:      os.Getenv("OPENAI_API_KEY"),
			AllowedOrigins: strings.Split(allowedOrigins, ","),
			GoogleOauth: &config.GoogleOauth{
				ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
				ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
				RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URI"),
			},
			WorkerRedisAddr: os.Getenv("REDIS_URL"),
		}

		s, err := server.NewServer(cfg, gormdb)
		if err != nil {
			// log.Fatal(fmt.Errorf("error creating server: %w", err))
			log.Error("error creating server", "error", err)
			os.Exit(1)
		}

		go func() {
			if err := s.ListenAndServe(); err != nil {
				// log.Fatal(fmt.Errorf("error running server: %w", err))
				log.Error("error running server", "error", err)
				os.Exit(1)
			}
		}()
	}

	if runWorker {
		concurrencyStr := os.Getenv("WORKER_CONCURRENCY")
		concurrency, err := strconv.Atoi(concurrencyStr)
		if err != nil {
			concurrency = 1
			log.Warn(
				"error parsing WORKER_CONCURRENCY, defaulting to 1", "error",
				fmt.Errorf("error parsing WORKER_CONCURRENCY: %q %w; defaulting to %d", concurrencyStr, err, concurrency),
			)
		}
		workerCfg := config.Worker{
			Addr:            os.Getenv("WORKER_ADDR"),
			Concurrency:     concurrency,
			OpenAIKey:       os.Getenv("OPENAI_API_KEY"),
			WorkerRedisAddr: os.Getenv("REDIS_URL"),
		}

		w, cancelWorker, err := worker.New(workerCfg, gormdb, rc)
		if err != nil {
			// log.Fatal(fmt.Errorf("error creating worker: %w", err))
			log.Error("error creating worker", "error", err)
			os.Exit(1)
		}
		workerCancel = cancelWorker

		// Run the worker in the background
		go func() {
			err := w.Run()
			if err != nil {
				// log.Fatal(fmt.Errorf("[worker] failed running worker: %w", err))
				log.Error("[worker] failed running worker", "error", err)
				os.Exit(1)
			}
			log.Info("Worker finished")
			close(completed)
		}()

		// Run health check in the background
		go func() {
			err := w.RunHealthServer()
			if err != nil {
				// log.Fatal(fmt.Errorf("error running health server: %w", err))
				log.Error("error running health server", "error", err)
				os.Exit(1)
			}
		}()

		go func() {

			convClient, err := conveyor.NewClient(&convconfig.Client{
				RedisURL:  os.Getenv("REDIS_URL"),
				Namespace: constants.DefaultJobNamespace,
				Logger:    log,
			})
			if err != nil {
				// log.Fatal(fmt.Errorf("error creating conveyor client: %w", err))
				log.Error("error creating conveyor client", "error", err)
				os.Exit(1)
			}
			time.Sleep(time.Second)
			// Ping it with a job
			_, err = convClient.Enqueue(context.Background(), &wire.Ping{})
			if err != nil {
				// log.Fatal(fmt.Errorf("error enqueuing job: %w", err))
				log.Error("error enqueuing job", "error", err)
				os.Exit(1)
			}
		}()
	}

	if runServer || runWorker {
		// Create a channel to listen for the termination signal
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		// Wait for the termination signal
		sig := <-signalChan
		log.Warn("Received signal: %s. Shutting down...", "signal", sig.String())
		if workerCancel != nil {
			workerCancel()
			log.Warn("Waiting for worker to finish...")
			<-completed
		}
	} else {
		log.Warn("Neither the server or the worker is running.")
	}

	log.Info("Goodbye! 👋")
}

func connectGorm() *gorm.DB {
	level := logger.Silent

	// if os.Getenv("ENV") != "development" || os.Getenv("LOG_SQL") != "false" {
	level = logger.Info
	// }

	newLogger := logger.New(
		log.Default(),
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  level,       // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)
	gormdb, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("error opening db: %w", err))
	}
	sqlDB, err := gormdb.DB()
	if err != nil {
		log.Fatal(fmt.Errorf("error getting db: %w", err))
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Fatal(fmt.Errorf("error pinging db: %w", err))
	}

	return gormdb
}

func connectToRedis(log *slog.Logger) *redis.Client {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		// log.Fatal(fmt.Errorf("error getting redis URL: %w", err))
		log.Error("error getting redis URL", "error", err)
		os.Exit(1)
	}
	rc := redis.NewClient(opt)
	_, err = rc.Ping(context.Background()).Result()
	if err != nil {
		// log.Fatal(fmt.Errorf("error pinging redis: %w", err))
		log.Error("error pinging redis", "error", err)
		os.Exit(1)
	}

	return rc
}

func connectToSentry(log *slog.Logger) {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	log.Info("Connecting to Sentry with environment", "env", env)
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		Environment:      env,
		EnableTracing:    true,
		AttachStacktrace: true,
		ServerName:       "reviso",
		TracesSampleRate: 1.0,
	}); err != nil {
		// log.Fatalf("Sentry initialization failed: %v\n", err)
		log.Error("Sentry initialization failed", "error", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println(`Usage: program [OPTIONS]
Options:
  --server      Run the server only
  --worker      Run the worker only
  -h            Show this help message

By default, without any options, both the server and worker will run.`)
}
