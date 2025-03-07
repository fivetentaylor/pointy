package env

import (
	"context"
	"log/slog"

	"github.com/charmbracelet/log"
	"github.com/jpoz/conveyor"
	"github.com/posthog/posthog-go"
	"github.com/redis/go-redis/v9"
	"github.com/sashabaranov/go-openai"
	"github.com/fivetentaylor/pointy/pkg/client"
	"github.com/fivetentaylor/pointy/pkg/query"
	"github.com/fivetentaylor/pointy/pkg/storage/dynamo"
	"github.com/fivetentaylor/pointy/pkg/storage/s3"
	"gorm.io/gorm"
)

func Attach(ctx context.Context,
	log *log.Logger,
	slog *slog.Logger,
	query *query.Query,
	oai *openai.Client,
	dynamo *dynamo.DB,
	ses client.SESInterface,
	rc *redis.Client,
	ph posthog.Client,
	bg *conveyor.Client,
	s3 *s3.S3,
	rawDB *gorm.DB,
) context.Context {
	ctx = QueryCtx(ctx, query)
	ctx = OpenAiCtx(ctx, oai)
	ctx = DynamoCtx(ctx, dynamo)
	ctx = SESCtx(ctx, ses)
	ctx = RedisCtx(ctx, rc)
	ctx = PosthogCtx(ctx, ph)
	ctx = BackgroundCtx(ctx, bg)
	ctx = LogCtx(ctx, log)
	ctx = SLogCtx(ctx, slog)
	ctx = S3Ctx(ctx, s3)
	ctx = RawDBCtx(ctx, rawDB)

	return ctx
}

func Copy(ctx context.Context) context.Context {
	newCtx := context.Background()

	newCtx = QueryCtx(newCtx, Query(ctx))
	newCtx = OpenAiCtx(newCtx, OpenAi(ctx))
	newCtx = DynamoCtx(newCtx, Dynamo(ctx))
	newCtx = RedisCtx(newCtx, Redis(ctx))
	newCtx = PosthogCtx(newCtx, Posthog(ctx))
	newCtx = SESCtx(newCtx, SES(ctx))
	newCtx = BackgroundCtx(newCtx, Background(ctx))
	newCtx = LogCtx(newCtx, Log(ctx))
	newCtx = SLogCtx(newCtx, SLog(ctx))
	newCtx = S3Ctx(newCtx, S3(ctx))
	newCtx = RawDBCtx(newCtx, RawDB(ctx))

	return newCtx
}
