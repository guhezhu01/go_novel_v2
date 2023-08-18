package middleware

import (
	"context"
	"fmt"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
)

func TracingClientMiddleWare(serviceName, id string, closer io.Closer) grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(
		func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			ctx = metadata.AppendToOutgoingContext(ctx, serviceName, id)
			defer func() {
				closer.Close()
			}()
			return invoker(ctx, method, req, reply, cc, opts...)
		},
		grpc_opentracing.UnaryClientInterceptor())
}

func GetRequestID(ctx context.Context, key string) (string, bool) {
	reqId, _ := ctx.Value(key).(string)
	if reqId == "" {
		md, _ := metadata.FromIncomingContext(ctx)
		reqIds := md.Get(key)
		if len(reqIds) == 0 {
			return reqId, false
		}
		reqId = reqIds[0]
		if reqId != "" {
			return reqId, true
		}
		return reqId, false
	}
	return reqId, true
}

func InitTracing(serviceName string) (opentracing.Tracer, io.Closer) {

	cfg := &config.Configuration{
		ServiceName: serviceName,
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
		panic(fmt.Sprintf("ERROR: cannot init-config Jaeger: %v\n", err))
	}

	return tracer, closer
}
