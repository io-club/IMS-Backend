package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	ioconfig "ims-server/pkg/config"
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

	// To implement simple permission checks, middleware needs to be run first
	svc.Use(ioginx.TimeMW(), ioginx.JwtAuthMW())

	svc.Any("/:service/:func", func(c *gin.Context) {
		service := c.Param("service")
		fn := c.Param("func")

		// Generate trace information
		c.Request.Header.Set("traceID", fmt.Sprintf("%d", rand.Int63()))
		c.Request.Header.Set("spanID", fmt.Sprintf("%d", rand.Int63()))

		// Service discovery
		hub := registery.GetServiceClient()
		servers := hub.GetServiceEndpointsWithCache(service)
		// TODO: Implement a more efficient load balancing algorithm
		// Load balancing (random method)
		idx := rand.Intn(len(servers))
		path := servers[idx]
		// Redirect service
		fn = strings.ToLower(fn)
		urlAddr := fmt.Sprintf("http://%s/%s", path, fn)
		iologger.Info("Redirect to " + urlAddr)

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

	ioginx.NewIOServer(svc).Run(addr, config.Name, nil)
}
