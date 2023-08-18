package main

import (
	"comment-service/discovery"
	initConfig "comment-service/init-config"
	"github.com/guhezhu01/go_novel_v2/model-tools/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	initConfig.Init()
	go func() {
		discovery.RegisterService()
	}()

	quit := make(chan os.Signal)
	// SIGINT 用户发送INTR(Ctrl+C)触发退出 SIGTERM结束程序
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("关闭用户服务！")
}
