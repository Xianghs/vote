package db

import (
	"github.com/go-redis/redis"
	_ "github.com/go-redis/redis"
	"log"
	"vote/vote/app/config"
)

var MyRedis *redis.Client

func InitRedis(config config.RedisConfig) {
	MyRedis = redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password,
		DB:       0,
	})

	if MyRedis == nil {
		log.Fatal("Redis连接失败！")
	}
}
