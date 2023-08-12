package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"user-service/config"
	"user-service/discovery"
	"user-service/internal/repository"
)

func main() {
	config.InitConfig()
	repository.InitDb()

	go func() {
		discovery.RegisterService()
	}()

	quit := make(chan os.Signal)
	// SIGINT 用户发送INTR(Ctrl+C)触发退出 SIGTERM结束程序
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("关闭用户服务！")

}
