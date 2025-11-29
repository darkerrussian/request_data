package requestdata

import (
	"context"
	"google.golang.org/grpc/metadata"
)

type contextKey int

const (
	ctxRequestDataKey = contextKey(iota)
	defaultTimeZone   = "Europe/Moscow"
)

func NewContextS(ctx context.Context, rd *RequestData) context.Context {
	return context.WithValue(metadata.NewOutgoingContext(ctx, rd.ToMD()), ctxRequestDataKey, rd)
}

func FromContext(ctx context.Context) *RequestData {
	if ctx == nil {
		return nil
	} else if val, ok := ctx.Value(ctxRequestDataKey).(*RequestData); ok {
		return val
	}
	return nil
}

func CopyCtx(dst context.Context, src context.Context) context.Context {
	return NewContextS(dst, FromContext(src))
}
