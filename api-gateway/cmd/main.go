package main

import (
	"api-gateway/config"
	"api-gateway/routes"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config.InitConfig()

	engine := routes.InitRoutes()

	go func() {
		engine.Run(viper.GetString("HttpPort"))
	}()

	quit := make(chan os.Signal)
	// SIGINT 用户发送INTR(Ctrl+C)触发退出 SIGTERM结束程序
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("关闭网关服务！")

}
