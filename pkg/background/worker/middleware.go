package worker

import (
	"context"

	"github.com/getsentry/sentry-go"
	cwire "github.com/jpoz/conveyor/wire"
)

func SentryMiddleware(ctx context.Context, job *cwire.Job, next func(context.Context) error) error {
	defer sentry.Recover()

	err := next(ctx)
	if err != nil {
		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("jobtype", job.Type)
			scope.SetContext("job", map[string]interface{}{
				"jid":     job.Uuid,
				"queue":   job.Queue,
				"jobtype": job.Type,
				"args":    job.Payload,
			})
		})
		sentry.CaptureException(err)
	}

	return err
}
