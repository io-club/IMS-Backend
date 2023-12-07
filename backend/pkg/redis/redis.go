package ioredis

import (
	"context"
	"github.com/redis/go-redis/v9"
	ioconfig "ims-server/pkg/config"
	iologger "ims-server/pkg/logger"
	"time"
)

func NewClient() *redis.Client {
	var rdb *redis.Client
	connectDoneCh := make(chan struct{})
	go func() {
		config := ioconfig.GetRedisConf()
		rdb = redis.NewClient(&redis.Options{
			Addr:     config.Addr,
			Password: config.Password,
			DB:       config.DB,
		})
		connectDoneCh <- struct{}{}
	}()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	select {
	case <-ctx.Done():
		iologger.Fatalf("redis client connect timeout")
	case <-connectDoneCh:
	}
	return rdb
}
