package initConfig

import (
	"comment-service/config"
	"comment-service/internal/repository"
	distributed_lock "github.com/guhezhu01/go_novel_v2/model-tools/distributed-lock"
	"github.com/guhezhu01/go_novel_v2/model-tools/log"
	"github.com/spf13/viper"
)

var DistributedLockConn *distributed_lock.RedisLock

func Init() {
	config.InitConfig()
	repository.InitDb()
	repository.InitCache()
	DistributedLockConn = distributed_lock.InitLock(viper.GetString("redis.Addr"), viper.GetString("redis.Password"))
	log.InitRpcLog(viper.GetString("LogPath"), viper.GetString("consul.Name"))
}
