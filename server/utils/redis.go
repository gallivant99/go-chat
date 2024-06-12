package utils

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

var Pool *redis.Pool

func InitRedisPool(maxIdle, maxActive int, idleTimeout time.Duration, host string) (*redis.Pool) {
	return &redis.Pool{
		// 初始化链接数量
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", host)
		},
	}
}