package booking

import (
	"context"
	"go-clean/src/lib/appcontext"
)

func (w *worker) initContext(ctx context.Context, reqId string) context.Context {
	ctx = appcontext.SetRequestID(ctx, reqId)
	return ctx
}
