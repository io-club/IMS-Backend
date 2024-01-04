package main

import (
	"context"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/redisstorage"
	"ims-server/microservices/work/internal/dal/model"
	"ims-server/microservices/work/internal/dal/repo"
	"ims-server/microservices/work/job"
	iocolly "ims-server/pkg/colly"
	ioconfig "ims-server/pkg/config"
	"ims-server/pkg/logger"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx := context.Background()
	config := ioconfig.GetServiceConf().Work
	serviceName := config.Name
	iolog.SetLogger(serviceName)

	go func() {
		// 删除创建时间为两小时以前的天气
		for {
			err := repo.NewWeatherRepo().DB().WithContext(ctx).Unscoped().Where("created_at < ?", time.Now().Add(-3*time.Hour)).Delete(&model.Weather{}).Error
			if err != nil {
				iolog.Info("delete weather failed: %v", err)
			}
		}
	}()

	prefix := "weather_storage"
	var storage *redisstorage.Storage
	go func() {
		collector := job.CreateWeatherCollector()
		storage = iocolly.SetRedisStorage(collector, prefix, 2*time.Hour)
		// 访问失败时打印日志
		collector.OnError(func(r *colly.Response, err error) {
			iolog.Warn(fmt.Sprintf("网页爬取失败: %s", err))
		})
		// 处理爬取数据
		collector.OnHTML("body", job.SaveWeatherData(ctx))
		// 找到带有 href 属性的标签 a
		collector.OnHTML(`a[href]`, func(e *colly.HTMLElement) {
			link := e.Attr("href")
			collector.Visit(e.Request.AbsoluteURL(link))
		})

		if err := collector.Visit("http://www.weather.com.cn/weather1d/101030500.shtml"); err != nil {
			iolog.Info("访问链接失败: %v", err)
		}
	}()
	// 优雅地终止程序
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	// 关闭前，清空 redis 中的 url 缓存
	storage.Clear()
	storage.Client.Close()
}
