package discovery

import (
	"article-service/internal/handler"
	"article-service/internal/service"
	"fmt"
	"github.com/hashicorp/consul/api"
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
	//初始化grpc对象
	grpcServer := grpc.NewServer()
	//注册服务
	service.RegisterArticleServiceServer(grpcServer, new(handler.ArticleService))
	//设置监听，指定ip/port
	addr := viper.GetString("consul.Address") + ":" + viper.GetString("consul.Port")
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err)
	}
	log.Println("启动文章服务！")

	//启动服务
	err = grpcServer.Serve(listen)
	if err != nil {
		fmt.Println(err)
	}
	defer listen.Close()

}
