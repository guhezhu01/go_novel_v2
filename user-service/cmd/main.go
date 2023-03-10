package main

import (
	"user-service/config"
	"user-service/discovery"
	"user-service/internal/repository"
)

func main() {
	config.InitConfig()
	repository.InitDb()
	discovery.RegisterService()

}
