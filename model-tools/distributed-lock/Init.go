package distributed_lock

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

type RedisLock struct {
	redisPool *redis.Pool
	id        string
	resource  string
	expire    time.Duration
}

// 对同一个数据加锁，但是只能一个进程解锁
func InitLock(addr, password string) (*RedisLock, bool) {
	pool := &redis.Pool{
		MaxIdle:     30,          // 最大连接数
		MaxActive:   100,         //最大活跃连接数，0代表无限
		IdleTimeout: time.Minute, //闲置连接的超时时间
		Dial: func() (redis.Conn, error) {

			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err

		},
	}
	lock := newRedisLock(pool, 10*time.Second)

	if lock == nil {
		return nil, false
	}
	return lock, true
}

func newRedisLock(pool *redis.Pool, expire time.Duration) *RedisLock {

	return &RedisLock{
		redisPool: pool,
		expire:    expire,
	}

}

// TryLock 加锁 argsId用于设置进程唯一id,防止其他进程解锁该进程操作的resource(默认不指定)
func (lock *RedisLock) TryLock(resource string, timeOut time.Duration, argsId ...string) bool {

	var id string
	if len(argsId) > 0 {
		id = argsId[0]
	}
	lock.id = id
	lock.resource = resource

	conn := lock.redisPool.Get()
	defer conn.Close()
	// 尝试获取锁
	now := time.Now()
	for {
		// set lock x EX 10 NX (EX设置过期时间,NX:只在键不存在的时候，才对键进行设置操作)
		// SET key value NX 等同于 SETNX key value
		result, err := redis.String(conn.Do("SET", resource, "1", "EX", int(lock.expire.Seconds()), "NX"))
		if err == nil && result == "OK" {
			log.Println("加锁成功!")
			return true
		}
		since := time.Since(now)
		// 判断是否超时
		if since >= timeOut {
			log.Println("尝试获取锁超时: ", err)
			return false
		}
	}

	return true

}

// Unlock 释放锁
func (lock *RedisLock) Unlock(argsId ...string) bool {
	if lock.id != "" {
		if len(argsId) > 0 {
			if argsId[0] != lock.id {
				fmt.Println("没有解锁权限, id = ", argsId[0])
				return false
			}
		} else {
			fmt.Println("解锁失败, id = ", argsId[0])
			return false
		}
	}

	conn := lock.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", lock.resource)

	if err != nil {
		fmt.Println("释放锁发生错误：", err)
		return false
	}
	fmt.Println("释放锁成功!")
	return true
}
