package asynctx

import (
	"context"
	"taxes-be/utils/logrus"

	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
)

func Linked(original context.Context, name string) context.Context {
	ctx := context.Background()
	// copy log fields
	ctx = logrus.Extend(ctx, logrus.Get(original)...)
	// link span
	oldSpan := trace.FromContext(original)
	if oldSpan != nil {
		var span *trace.Span
		ctx, span = trace.StartSpan(ctx, name)

		sc := oldSpan.SpanContext()
		span.AddLink(trace.Link{
			TraceID: sc.TraceID,
			SpanID:  sc.SpanID,
		})
	}
	// metric's tags
	ctx = tag.NewContext(ctx, tag.FromContext(original))

	return ctx
}
