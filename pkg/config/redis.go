package ioconfig

import (
	"log"
	"sync"
)

var (
	redisConf *RedisConf
	redisOnce sync.Once
)

type RedisConf struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func GetRedisConf() *RedisConf {
	redisOnce.Do(func() {
		if redisConf == nil {
			if err := V.UnmarshalKey("redis", &redisConf); redisConf == nil || err != nil {
				log.Panicf("unmarshal conf failed, err: %s \n", err)
			}
		}
	})
	return redisConf
}
