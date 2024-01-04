package iocolly

import (
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/redisstorage"
	ioconfig "ims-server/pkg/config"
	"time"
)

// SetRedisStorage 该函数返回一个 redis 的 client，应显式关闭该 client
func SetRedisStorage(c *colly.Collector, prefix string, expires time.Duration) *redisstorage.Storage {
	// 将 Cookie 和访问过的 URL 存储到 redis 中，并设置两小时的超时时间
	config := ioconfig.GetRedisConf()
	storage := &redisstorage.Storage{
		Address:  config.Addr,
		Password: config.Password,
		DB:       config.DB,
		Prefix:   prefix,
		Expires:  expires,
	}
	// 设置 Storage
	if err := c.SetStorage(storage); err != nil {
		panic(err)
	}
	return storage
}
