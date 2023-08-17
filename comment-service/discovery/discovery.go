package discovery

import (
	"comment-service/internal/handler"
	"comment-service/internal/service"
	"comment-service/middleware"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func RegisterService() {
	consulConfig := api.DefaultConfig()
	//创建consul对象
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		fmt.Println(err)
	}
	//注册的配置信息
	reg := api.AgentServiceRegistration{
		ID:      viper.GetString("consul.ID"),
		Tags:    viper.GetStringSlice("consul.Tags"),
		Name:    viper.GetString("consul.Name"),
		Address: viper.GetString("consul.Address"),
		Port:    viper.GetInt("consul.Port"),
		Check: &api.AgentServiceCheck{
			CheckID:  viper.GetString("consul.Check.CheckID"),
			TCP:      viper.GetString("consul.Check.TCP"),
			Timeout:  viper.GetString("consul.Check.Timeout"),
			Interval: viper.GetString("consul.Check.Interval"),
		},
	}
	//注册grpc服务到consul上

	err = consulClient.Agent().ServiceRegister(&reg)
	if err != nil {
		log.Println(err)
		os.Exit(0)
	}
	tracer, closer := middleware.InitTracing("comment-service")
	defer func() {
		closer.Close()
	}()
	opentracing.SetGlobalTracer(tracer)

	//初始化grpc对象
	grpcServer := grpc.NewServer(middleware.Tracing("comment-service"), middleware.AuthCheckToken("comment-service"))
	//注册服务
	service.RegisterCommentServiceServer(grpcServer, new(handler.CommentsService))
	//设置监听，指定ip/port
	addr := viper.GetString("consul.Address") + ":" + viper.GetString("consul.Port")

	//logger.WithContext(context.Background(), "comment-service").Infof("comment-service")
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err)
	}
	log.Println("启动评论服务！")

	//启动服务
	err = grpcServer.Serve(listen)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		err := listen.Close()
		if err != nil {
			fmt.Println("评论服务关闭失败！")
		}
	}()

}
