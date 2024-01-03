package job

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gocolly/colly/v2"
	"ims-server/microservices/work/internal/dal/model"
	"ims-server/microservices/work/internal/dal/repo"
	iocolly "ims-server/pkg/colly"
	iolog "ims-server/pkg/logger"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Weather struct {
	Od Od `json:"od"`
}

type Od struct {
	Od0 string `json:"od0"` // 日期
	Od1 string `json:"od1"` // 地点
	Od2 []Od2  `json:"od2"` // 天气数据
}

type Od2 struct {
	Od21 string `json:"od21"` // 时间
	Od22 string `json:"od22"` // 温度
	Od23 string `json:"od23"`
	Od24 string `json:"od24"` // 风向
	Od25 string `json:"od25"` // 风力
	Od26 string `json:"od26"` // 降水量
	Od27 string `json:"od27"` // 湿度
	Od28 string `json:"od28"`
}

// RemoveNewlinesAndTabs 去除 s 中的所有换行符和制表符
func RemoveNewlinesAndTabs(s string) string {
	blankReg := regexp.MustCompile(`[\n\t]`)
	str := blankReg.ReplaceAllString(s, "")
	return str
}

// SaveWeatherData 爬取并解析、存储天气数据
func SaveWeatherData(ctx context.Context) func(e *colly.HTMLElement) {
	return func(e *colly.HTMLElement) {
		dom := e.DOM

		location := dom.Find("div.crumbs.fl").Text()
		location = RemoveNewlinesAndTabs(location)
		// 去除空格并将>替换为 -
		location = strings.ReplaceAll(location, " ", "")
		location = strings.ReplaceAll(location, ">", "-")

		data := Weather{}
		dataJson := dom.Find("div.left-div script:nth-child(1)").Text()
		dataJson = strings.ReplaceAll(dataJson, " ", "")
		parts := strings.Split(dataJson, "=")
		if len(parts) >= 2 {
			dataJson = parts[1]
		} else {
			iolog.Info("获取到了无法识别的数据，请检查爬虫正确性")
			return
		}
		dataJson = RemoveNewlinesAndTabs(dataJson)
		// TODO：会报错，但显示忽略错误后却能正常运行，后期应检查
		_ = sonic.Unmarshal([]byte(dataJson), &data)

		var weathers []*model.Weather
		for i, od2 := range data.Od.Od2 {
			temperature, err := strconv.ParseFloat(od2.Od22, 32)
			if err != nil {
				temperature = 0
			}
			windForce, err := strconv.Atoi(od2.Od25)
			if err != nil {
				temperature = 0
			}
			precipitation, err := strconv.Atoi(od2.Od26)
			if err != nil {
				temperature = 0
			}
			humidity, err := strconv.Atoi(od2.Od27)
			if err != nil {
				temperature = 0
			}

			weather := model.Weather{
				DateTime: data.Od.Od0,
				Location: location[7:], // 去除"全国-"前缀

				Hour:          od2.Od21,
				Temperature:   int8(temperature),
				WindDirection: od2.Od24,
				WindForce:     int8(windForce),
				Precipitation: int8(precipitation),
				Humidity:      int8(humidity),
			}
			weathers = append(weathers, &weather)
			// 只抓取 12 小时的天气数据
			if i == 12 {
				break
			}
		}
		err := repo.NewWeatherRepo().MCreate(ctx, weathers)
		if err != nil {
			iolog.Warn("插入天气数据失败: %v", err)
			return
		}
	}
}

func CrawlWeather(ctx context.Context) {
	collector := colly.NewCollector(
		// 限制域名
		colly.AllowedDomains("www.weather.com.cn"),
		// 用正则限制爬取范围
		colly.URLFilters(
			regexp.MustCompile(`.*/weather1d/101.+`),
		),
		colly.DisallowedURLFilters(
			regexp.MustCompile(`.*A.shtml$`), // 去除景区数据
		),
		iocolly.RotateUserAgent(),
		// 限制并发数与间隔，避免被线下重拳出击
		iocolly.SetRateLimit("*", 3, rand.Intn(4)+1),
		// TODO：增加随机 IP，进一步防止被重拳出击
	)
	client := iocolly.SetRedisStorage(collector, "weather_storage", 2*time.Hour)
	defer client.Close()

	// 访问失败时打印日志
	collector.OnError(func(r *colly.Response, err error) {
		iolog.Warn(fmt.Sprintf("网页爬取失败: %s", err))
	})

	collector.OnHTML("body", SaveWeatherData(ctx))
	// 找到带有 href 属性的标签 a
	collector.OnHTML(`a[href]`, func(e *colly.HTMLElement) {
		link := e.Attr("href")
		collector.Visit(e.Request.AbsoluteURL(link))
	})

	if err := collector.Visit("http://www.weather.com.cn/weather1d/101030500.shtml"); err != nil {
		iolog.Info("访问链接失败: %v", err)
	}
	return
}
