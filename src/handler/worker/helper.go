package worker

import (
	"context"
	"go-clean/src/lib/appcontext"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func (w *worker) initContext(ctx context.Context, reqId string, traceContext map[string]string) (context.Context, error) {
	if traceContext != nil {
		propagator := otel.GetTextMapPropagator()
		carrier := propagation.MapCarrier(traceContext)
		ctx = propagator.Extract(context.Background(), carrier)
	}

	ctx = appcontext.SetRequestID(ctx, reqId)

	return ctx, nil
}
