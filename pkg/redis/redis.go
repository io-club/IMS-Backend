package ioredis

import (
	"github.com/redis/go-redis/v9"
	ioconfig "ims-server/pkg/config"
)

func NewClient() *redis.Client {
	config := ioconfig.GetRedisConf()
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})
	return rdb
}
