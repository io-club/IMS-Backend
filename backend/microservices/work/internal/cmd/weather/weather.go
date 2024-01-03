package main

import (
	"context"
	"ims-server/microservices/work/internal/dal/model"
	"ims-server/microservices/work/internal/dal/repo"
	"ims-server/microservices/work/job"
	ioconfig "ims-server/pkg/config"
	"ims-server/pkg/logger"
	"time"
)

func main() {
	ctx := context.Background()
	config := ioconfig.GetServiceConf().Work
	serviceName := config.Name
	iolog.SetLogger(serviceName)

	// 删除创建时间为两小时以前的天气
	go func() {
		for {
			err := repo.NewWeatherRepo().DB().WithContext(ctx).Unscoped().Where("created_at < ?", time.Now().Add(-2*time.Hour)).Delete(&model.Weather{}).Error
			if err != nil {
				iolog.Info("delete weather failed: %v", err)
			}
		}
	}()
	job.CrawlWeather(ctx)
	// TODO:关闭时，应清理 redis 中的缓存
}
