package db

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type RedisConf struct {
	Addr    string
	Pwd     string
	Db      int
	Size    int
	Timeout int
}

func InitRedis(conf *RedisConf) *redis.Client {
	redisConnection := redis.NewClient(&redis.Options{
		Addr:        conf.Addr,
		Password:    conf.Pwd,
		DB:          conf.Db,
		PoolSize:    conf.Size,
		PoolTimeout: time.Duration(conf.Timeout),
	})
	if _, err := redisConnection.Ping().Result(); err != nil {
		panic(err)
	} else {
		fmt.Printf("redis init success, addr:%s", conf.Addr)
	}
	return redisConnection
}
