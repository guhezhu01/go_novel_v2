package discovery

import (
	"api-gateway/internal/service"
	"api-gateway/middleware"
	"context"
	"fmt"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/hashicorp/consul/api"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"strconv"
)

func GetService(ctx context.Context, serviceName, tag string) interface{} {
	//logger.WithContext(ctx).Infof("调用的服务:%s", serviceName)
	//初始化consul配置
	config := api.DefaultConfig()
	config.Address = viper.GetString("consul.Address") + ":" + viper.GetString("consul.Port")
	//创建consul对象
	consulClient, err01 := api.NewClient(config)

	if err01 != nil {
		fmt.Println(err01)
	}
	//服务发现，从consul上获取健康的服务
	services, _, err := consulClient.Health().Service(serviceName, tag, true, nil)

	//使用从服务发现获取到服务的ip/port
	addr := services[0].Service.Address + ":" + strconv.Itoa(services[0].Service.Port)

	tracer, closer := middleware.InitTracing(serviceName)
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	//连接服务
	grpcConn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithChainUnaryInterceptor(
		func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			id, ok := middleware.GetRequestID(ctx, serviceName)
			if !ok {
				return nil
			}
			ctx = metadata.AppendToOutgoingContext(ctx, "x-request-id", id)
			return invoker(ctx, method, req, reply, cc, opts...)
		},
		grpc_opentracing.UnaryClientInterceptor(),
	))
	if err != nil {
		fmt.Println("连接错误:", err)
	}

	//初始化grpc客户端
	var grpcClient interface{}
	switch serviceName {
	case "user service":
		grpcClient = service.NewUserServiceClient(grpcConn)
	case "comment service":
		grpcClient = service.NewCommentServiceClient(grpcConn)
	}
	return grpcClient
}
