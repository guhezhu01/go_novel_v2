package distributed_lock

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

type RedisLock struct {
	redisPool *redis.Pool
	resource  string
	expire    time.Duration
}

func InitLock(resource string) *RedisLock {
	pool := &redis.Pool{
		MaxIdle:     30,          // 最大连接数
		MaxActive:   100,         //最大活跃连接数，0代表无限
		IdleTimeout: time.Minute, //闲置连接的超时时间
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
	lock := NewRedisLock(pool, resource, 10*time.Second)
	return lock
}

func NewRedisLock(pool *redis.Pool, resource string, expire time.Duration) *RedisLock {

	return &RedisLock{
		redisPool: pool,
		resource:  resource,
		expire:    expire,
	}

}

func (lock *RedisLock) TryLock() bool {

	conn := lock.redisPool.Get()

	defer conn.Close()

	// 尝试获取锁

	result, err := redis.String(conn.Do("SET", lock.resource, "1", "EX", int(lock.expire.Seconds()), "NX"))

	if err != nil {

		fmt.Println("尝试获取锁发生错误：", err)

		return false

	}

	return result == "OK"

}

func (lock *RedisLock) Unlock() bool {

	conn := lock.redisPool.Get()

	defer conn.Close()

	_, err := conn.Do("DEL", lock.resource)

	if err != nil {
		fmt.Println("释放锁发生错误：", err)
		return false
	}
	return true
}
