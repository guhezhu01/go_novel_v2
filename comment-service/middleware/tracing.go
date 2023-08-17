package middleware

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/guhezhu01/go_novel_v2/model-tools/log"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
)

func Tracing(key string) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(withRequestID(key), grpc_opentracing.UnaryServerInterceptor())
}

func withRequestID(key string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, _ := metadata.FromIncomingContext(ctx)
		reqId := md.Get(key)
		log.Println(reqId)
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

func InitTracing(service string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: "123.249.88.132:6831",
			LogSpans:           true,
		},
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		log.Println(err)
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}

	return tracer, closer
}
