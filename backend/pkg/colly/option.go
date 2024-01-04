package iocolly

import (
	"bytes"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/gocolly/colly/v2/proxy"
	ioconfig "ims-server/pkg/config"
	"ims-server/pkg/util"
	"log"
	"net"
	"net/http"
	"time"
)

// 打印抓取到的完整 html 页面。Debug 时使用，先确保能拿到 html 页面再去写解析 html 的代码
func PrintResponse() colly.CollectorOption {
	return func(c *colly.Collector) {
		c.OnResponse(func(r *colly.Response) {
			log.Printf("%s\n", bytes.Replace(r.Body, []byte("\n"), nil, -1))
		})
	}
}

// 使用 socks5 代理
func SetProxy() colly.CollectorOption {
	return func(c *colly.Collector) {
		proxies := make([]string, 0, 300)
		for _, ip := range util.ReadAllLines(ioconfig.RootPath + "/conf/socks5.txt") {
			fmt.Printf("ip: %s\n", "socks5://"+ip)
			proxies = append(proxies, "socks5://"+ip)
		}
		roundProxy, err := proxy.RoundRobinProxySwitcher(proxies...)
		if err != nil {
			log.Fatal(err)
		}
		c.SetProxyFunc(roundProxy)
	}
}

// 并发&限速
func SetRateLimit(DomainGlob string, Parallelism int, RandomDelay int) colly.CollectorOption {
	return func(c *colly.Collector) {
		err := c.Limit(&colly.LimitRule{
			DomainGlob:  DomainGlob,
			Parallelism: Parallelism, //并发
			// Delay:       5 * time.Second,	//定长休息
			RandomDelay: time.Duration(RandomDelay) * time.Second, //随机休息
		})
		if err != nil {
			panic(err)
		}
	}
}

// 设置超时
func SetTimeout(DialTimeout, TLSHandshakeTimeout time.Duration) colly.CollectorOption {
	return func(c *colly.Collector) {
		c.WithTransport(&http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   DialTimeout,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   TLSHandshakeTimeout,
			ExpectContinueTimeout: 1 * time.Second,
		})
	}
}

// RotateUserAgent 模拟浏览器生成 User-Agent，并随机切换。注意：colly 自带的 User-Agent 数量是有限的
func RotateUserAgent() colly.CollectorOption {
	return func(c *colly.Collector) {
		extensions.RandomUserAgent(c) //生成一些逼真的 Chrome/Firefox/Opera User-Agent
	}
}
