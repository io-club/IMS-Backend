package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	ioconfig "ims-server/pkg/config"
	ioetcd "ims-server/pkg/etcd"
	ioginx "ims-server/pkg/ginx"
	iologger "ims-server/pkg/logger"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	config := ioconfig.GetServiceConf().Nms
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)

	svc := gin.Default()
	svc.Any("/:service/:func", func(c *gin.Context) {
		service := c.Param("service")
		fn := c.Param("func")

		publicMap := GetPublicRouteMap(NmsRoutesMap)
		// 禁止访问私有服务
		if _, ok := publicMap[service][strings.ToLower(fn)]; !ok {
			iologger.Warn("用户异常访问内部服务, service: %s, fn: %s", service, fn)
			c.String(http.StatusForbidden, "")
			return
		}

		// 生成 trace 信息
		c.Request.Header.Set("traceID", fmt.Sprintf("%d", rand.Int63()))
		c.Request.Header.Set("spanID", fmt.Sprintf("%d", rand.Int63()))

		// 服务发现
		hub := ioetcd.GetServiceHub()
		servers := hub.GetServiceEndpointsWithCache(ioetcd.UserService)
		// TODO：实现更高效的负载均衡算法
		// 负载均衡（随机法）
		idx := rand.Intn(len(servers))
		addr := servers[idx]
		// 重定向服务
		fn = strings.ToLower(fn)
		urlAddr := fmt.Sprintf("http://%s/%s", addr, fn)
		iologger.Info("redirect to " + urlAddr)

		remote, err := url.Parse(urlAddr)
		if err != nil {
			panic(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
		// Define the director func
		// This is a good place to log, for example
		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = remote.Host
			req.Method = c.Request.Method
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = remote.Path
		}

		proxy.ModifyResponse = func(response *http.Response) error {
			return nil
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	})

	ioginx.NewIOServer(svc).Run(addr, config.Name, nil, "")
}
