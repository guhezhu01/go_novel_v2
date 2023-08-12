package main

import (
	"comment-service/config"
	"comment-service/discovery"
	"comment-service/internal/repository"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	config.InitConfig()
	repository.InitDb()
	repository.InitCache()

	go func() {
		discovery.RegisterService()
	}()

	quit := make(chan os.Signal)
	// SIGINT 用户发送INTR(Ctrl+C)触发退出 SIGTERM结束程序
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("关闭用户服务！")
}
