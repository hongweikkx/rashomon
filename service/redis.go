package service

import (
	"runtime"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hongweikkx/rashomon/conf"
)

var RedisClient *redis.Client

func init() {
	opt := redis.Options{
		Addr:         conf.AppConfig.RedisHost,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		PoolSize:     runtime.NumCPU() / 2,
		MinIdleConns: 1,
	}
	RedisClient = redis.NewClient(&opt)
}
