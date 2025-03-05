package testutils

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/jpoz/conveyor"
	"github.com/jpoz/conveyor/config"
	"github.com/redis/go-redis/v9"
	"github.com/sashabaranov/go-openai"
	"github.com/teamreviso/code/pkg/client"
	"github.com/teamreviso/code/pkg/constants"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/models"
	"github.com/teamreviso/code/pkg/query"
	"github.com/teamreviso/code/pkg/storage/dynamo"
	"github.com/teamreviso/code/pkg/storage/s3"
	"github.com/teamreviso/code/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func WithUserClaim(user *models.UserClaims) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return env.UserClaimCtx(ctx, user)
	}
}

func WithUserClaimForUser(user *models.User) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return env.UserClaimCtx(ctx, &models.UserClaims{
			Id:    user.ID,
			Email: user.Email,
			Admin: user.Admin,
		})
	}
}

func TestContext(with ...func(context.Context) context.Context) context.Context {
	ctx := context.Background()

	for _, f := range with {
		ctx = f(ctx)
	}

	level, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err == nil {
		log.SetLevel(level)
	}

	// Gorm
	newLogger := logger.New(
		log.Default(),
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			// LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false, // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,  // Don't include params in the SQL log
			Colorful:                  true,  // Disable color
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

	// Redis
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal(fmt.Errorf("error getting redis URL: %w", err))
	}
	rc := redis.NewClient(opt)
	_, err = rc.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(fmt.Errorf("error pinging redis: %w", err))
	}

	q := query.Use(gormdb)

	openAiClient := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	dynamo, err := dynamo.NewDB()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect to dynamodb: %w", err))
	}

	bg, err := conveyor.NewClient(&config.Client{
		RedisURL:  os.Getenv("REDIS_URL"),
		Namespace: constants.DefaultJobNamespace,
		Logger:    slog.Default(),
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect to background worker: %w", err))
	}

	emailClient, err := client.NewSESFromEnv(bg)
	if err != nil {
		log.Fatal(fmt.Errorf("error creating email client: %w", err))
	}

	s3, err := s3.NewS3()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect to s3: %w", err))
	}

	ctx = env.Attach(
		ctx,
		log.Default(),
		utils.NewSlogFromEnv(),
		q,
		openAiClient,
		dynamo,
		emailClient,
		rc,
		nil, // posthog
		bg,
		s3,
		gormdb,
	)

	return ctx

}
