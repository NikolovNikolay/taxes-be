package cron

import (
	"context"
	"taxes-be/internal/atleastonce"
	"time"

	"github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

func RunInternalTimer(alod *atleastonce.Doer, limit int) {
	for range time.Tick(1 * time.Minute) {
		ctx := context.Background()
		ctx, span := trace.StartSpan(ctx, "cron.RunInternalTimer.Minute")
		everyMinute(ctx, alod, limit)
		span.End()
	}
}

func everyMinute(ctx context.Context, alod *atleastonce.Doer, limit int) {
	logrus.WithContext(ctx).Debug("periodic tasks started")
	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	lCtx, span := trace.StartSpan(ctx, "atleastonce.Doer.TryAll")
	defer span.End()

	alod.TryAll(lCtx, limit)

	logrus.WithContext(ctx).Debug("periodic tasks finished")
}
