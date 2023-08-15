package middleware

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func Tracing(key string) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(withRequestID(key))
}

func withRequestID(key string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, _ := metadata.FromIncomingContext(ctx)
		reqId := md.Get(key)
		if len(reqId) > 0 {
			ctx = NewWithContext(ctx, reqId[0], key)
		}
		return handler(ctx, req)
	}
}

func NewWithContext(ctx context.Context, id string, key string) context.Context {
	if id == "" {
		id = uuid.New().String()
	}
	return context.WithValue(ctx, key, id)
}

func GetRequestID(ctx context.Context, key string) (string, bool) {
	reqId, _ := ctx.Value(key).(string)
	if reqId == "" {
		return reqId, false
	}
	return reqId, true
}
