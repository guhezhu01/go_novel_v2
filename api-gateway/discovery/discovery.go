package discovery

import (
	"api-gateway/internal/service"
	"api-gateway/middleware"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"strconv"
)

func GetService(serviceName, tag string) interface{} {

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

	tracer, closer := middleware.InitTracing(viper.GetString("ServiceName"))
	opentracing.SetGlobalTracer(tracer)
	//连接服务
	id := uuid.New().String()
	grpcConn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		middleware.TracingClientMiddleWare(serviceName, id, closer),
		middleware.AuthMiddleWare(viper.GetString("AuthKey"), id),
	)
	if err != nil {
		fmt.Println("连接错误:", err)
	}

	//初始化grpc客户端
	var grpcClient interface{}
	switch serviceName {
	case "user-service":
		grpcClient = service.NewUserServiceClient(grpcConn)
	case "comment-service":
		grpcClient = service.NewCommentServiceClient(grpcConn)
	}
	return grpcClient
}
