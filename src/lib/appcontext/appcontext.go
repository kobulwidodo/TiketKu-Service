package appcontext

import (
	"context"
)

type contextKey string

const (
	requestId contextKey = "RequestId"
)

func GetRequestID(ctx context.Context) string {
	rqid, ok := ctx.Value(requestId).(string)
	if !ok {
		return ""
	}

	return rqid
}

func SetRequestID(ctx context.Context, rqid string) context.Context {
	return context.WithValue(ctx, requestId, rqid)
}
