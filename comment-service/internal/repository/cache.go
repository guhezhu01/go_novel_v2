package repository

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"log"
	"time"
)

var redisCache *redis.Client

func InitCache() {
	redisCache = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.Addr"),
		Password:     viper.GetString("redis.Password"),
		DB:           viper.GetInt("redis.DB"),                                        // redis一共16个库，指定其中一个库即可
		PoolSize:     viper.GetInt("redis.PoolSize"),                                  //最大连接数 默认cpu数*10
		ReadTimeout:  time.Duration(viper.GetInt("redis.ReadTimeout")) * time.Second,  //取超时间 单位秒 默认值为3秒
		WriteTimeout: time.Duration(viper.GetInt("redis.WriteTimeout")) * time.Second, // 写入超时时间 单位秒 默认与读超时时间一致
	})
	_, err := redisCache.Ping().Result()
	if err != nil {
		log.Fatalln(fmt.Sprintf("Redis连接错误:%s", err.Error()))
	} else {
		log.Println("Redis连接成功！")
	}

	return

}
