package discovery

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
	"user-service/internal/handler"
	"user-service/internal/service"
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

	fmt.Println(reg.Check.CheckID)
	consulClient.Agent().ServiceRegister(&reg)

	//初始化grpc对象
	grpcServer := grpc.NewServer()
	//注册服务
	service.RegisterUserServiceServer(grpcServer, new(handler.UserService))
	//设置监听，指定ip/port
	addr := viper.GetString("consul.Address") + ":" + viper.GetString("consul.Port")
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("开始监听")
	//启动服务
	grpcServer.Serve(listen)
	fmt.Println("监听结束")
	defer listen.Close()
}
