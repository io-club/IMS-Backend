package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ioconfig "ims-server/pkg/config"
	ioconsts "ims-server/pkg/consts"
	egoerror "ims-server/pkg/error"
	ioginx "ims-server/pkg/ginx"
	iologger "ims-server/pkg/logger"
	"ims-server/pkg/registery"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	config := ioconfig.GetServiceConf().Nms
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)

	iologger.SetLogger(config.Name)

	svc := gin.Default()
	// 启用代理服务
	svc.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		// AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin == "https://github.com"
		// },
		// MaxAge: 12 * time.Hour,
	}))

	// 为了实现简单权限检查，得先运行中间件
	svc.Use(ioginx.TimeMW(), ioginx.JwtAuthMW())

	svc.Any("/:service/:func", func(c *gin.Context) {
		service := c.Param("service")
		fn := c.Param("func")

		publicMap := GetPublicRouteMap(NmsRoutesMap)
		// 禁止访问私有服务
		route, ok := publicMap[service][strings.ToLower(fn)]
		if !ok {
			iologger.Warn("用户异常访问内部服务, service: %s, fn: %s", service, fn)
			c.String(http.StatusMethodNotAllowed, "")
			return
		}
		// 检查权限
		if route.Permission != nil && !route.Permission.IsEmpty() {
			utype, ok := c.Get("utype")
			if !ok {
				c.JSON(http.StatusUnauthorized, ioginx.NewErr(c, egoerror.ErrUnauthorized))
				return
			}
			userType := ioconsts.UserType(utype.(string))
			if _, ok := route.Permission[userType]; !ok {
				iologger.Debug("用户权限不足, service: %s, fn: %s", service, fn)
				c.JSON(http.StatusUnauthorized, ioginx.NewErr(c, egoerror.ErrNotPermitted))
				return
			}
		}
		// 生成 trace 信息
		c.Request.Header.Set("traceID", fmt.Sprintf("%d", rand.Int63()))
		c.Request.Header.Set("spanID", fmt.Sprintf("%d", rand.Int63()))

		// 服务发现
		hub := registery.GetServiceClient()
		servers := hub.GetServiceEndpointsWithCache(registery.UserService)
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

		// TODO:2023/11/29 21:21:31 http: proxy error: EOF   POST     "/user/SendVerification"
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

	ioginx.NewIOServer(svc).Run(addr, config.Name)
}
