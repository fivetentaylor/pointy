package env

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"runtime/debug"

	"github.com/charmbracelet/log"
	"github.com/jpoz/conveyor"
	"github.com/posthog/posthog-go"
	"github.com/redis/go-redis/v9"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"

	"github.com/teamreviso/code/pkg/client"
	"github.com/teamreviso/code/pkg/models"
	"github.com/teamreviso/code/pkg/query"
	"github.com/teamreviso/code/pkg/storage/dynamo"
	"github.com/teamreviso/code/pkg/storage/s3"
)

type ContextKey[T any] struct{}

var openAiKey = ContextKey[*openai.Client]{}
var OpenAiCtx = WithValueFunc(openAiKey)
var OpenAi = ValueFunc(openAiKey)

var dynamoKey = ContextKey[*dynamo.DB]{}
var DynamoCtx = WithValueFunc(dynamoKey)
var Dynamo = ValueFunc(dynamoKey)

var queryKey = ContextKey[*query.Query]{}
var QueryCtx = WithValueFunc(queryKey)
var Query = ValueFunc(queryKey)

var rawDBKey = ContextKey[*gorm.DB]{}
var RawDBCtx = WithValueFunc(rawDBKey)
var RawDB = ValueFunc(rawDBKey)

var userClaimKey = ContextKey[*models.UserClaims]{}
var UserClaimCtx = WithValueFunc(userClaimKey)
var UserClaim = ValueFuncWithError(userClaimKey)

var redisKey = ContextKey[*redis.Client]{}
var RedisCtx = WithValueFunc(redisKey)
var Redis = ValueFunc(redisKey)

var posthogKey = ContextKey[posthog.Client]{}
var PosthogCtx = WithValueFunc(posthogKey)
var Posthog = ValueFunc(posthogKey)

var s3Key = ContextKey[*s3.S3]{}
var S3Ctx = WithValueFunc(s3Key)
var S3 = ValueFunc(s3Key)

var sesKey = ContextKey[client.SESInterface]{}
var SESCtx = WithValueFunc(sesKey)
var SES = ValueFunc(sesKey)

var bgKey = ContextKey[*conveyor.Client]{}
var BackgroundCtx = WithValueFunc(bgKey)
var Background = ValueFunc(bgKey)

var logKey = ContextKey[*log.Logger]{}
var LogCtx = WithValueFunc(logKey)
var Log = ValueFunc(logKey)

var slogKey = ContextKey[*slog.Logger]{}
var SLogCtx = WithValueFunc(slogKey)
var SLog = ValueFunc(slogKey)

func WithValueFunc[T any](key ContextKey[T]) func(context.Context, T) context.Context {
	return func(ctx context.Context, value T) context.Context {
		return context.WithValue(ctx, key, value)
	}
}

func ValueFunc[T any](key ContextKey[T]) func(context.Context) T {
	return func(ctx context.Context) T {
		val, ok := ctx.Value(key).(T)
		if !ok {
			slog.Error("[PANIC] value not found in context for the given key", slog.Any("key", key), slog.String("stack", string(debug.Stack())))
			panic(fmt.Errorf("value not found in context for the given key: %+v with type: %T", key, key))
		}
		return val
	}
}

func ValueFuncWithError[T any](key ContextKey[T]) func(context.Context) (T, error) {
	return func(ctx context.Context) (T, error) {
		val, ok := ctx.Value(key).(T)
		if !ok {
			val = *new(T)
			reflectVal := reflect.ValueOf(*new(T))
			return *new(T), fmt.Errorf("value not found in context for the given key: %+v with type: %T", key, reflectVal.Type())
		}
		return val, nil
	}
}
