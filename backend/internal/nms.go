package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	deviceroute "ims-server/internal/device/route"
	userroute "ims-server/internal/user/route"
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

var NmsRoutesMap = map[string][]ioginx.Route{
	ioconfig.GetServiceConf().User.Name:   userroute.Routes,
	ioconfig.GetServiceConf().Device.Name: deviceroute.Routes,
}

func GetPublicRouteMap(fullServiceFuncList map[string][]ioginx.Route) map[string]map[string]ioginx.Route {
	publicRouteMap := map[string]map[string]ioginx.Route{}
	for service, routes := range fullServiceFuncList {
		publicRoutes := []ioginx.Route{}
		for _, route := range routes {
			if route.Private {
				continue
			}
			publicRoutes = append(publicRoutes, route)
		}

		for _, route := range publicRoutes {
			if err := ioginx.ParseRoute(&route); err != nil {
				panic(err)
			}
		}

		indexMap := map[string]ioginx.Route{}
		for _, route := range publicRoutes {
			indexMap[strings.ToLower(route.FuncName)] = route
		}
		publicRouteMap[service] = indexMap
	}
	return publicRouteMap
}

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

		publicMap := GetPublicRouteMap(NmsRoutesMap)
		// Deny access to private services
		route, ok := publicMap[service][strings.ToLower(fn)]
		if !ok {
			iologger.Warn("User accessed internal service abnormally, service: %s, fn: %s", service, fn)
			c.String(http.StatusMethodNotAllowed, "")
			return
		}
		// Check permissions
		if route.Permission != nil && !route.Permission.IsEmpty() {
			utype, ok := c.Get("utype")
			if !ok {
				c.JSON(http.StatusUnauthorized, ioginx.NewErr(c, egoerror.ErrUnauthorized))
				return
			}
			userType := ioconsts.UserType(utype.(string))
			if _, ok := route.Permission[userType]; !ok {
				iologger.Debug("Insufficient user permissions, service: %s, fn: %s", service, fn)
				c.JSON(http.StatusUnauthorized, ioginx.NewErr(c, egoerror.ErrNotPermitted))
				return
			}
		}
		// Generate trace information
		c.Request.Header.Set("traceID", fmt.Sprintf("%d", rand.Int63()))
		c.Request.Header.Set("spanID", fmt.Sprintf("%d", rand.Int63()))

		// Service discovery
		hub := registery.GetServiceClient()
		servers := hub.GetServiceEndpointsWithCache(registery.UserService)
		// TODO: Implement a more efficient load balancing algorithm
		// Load balancing (random method)
		idx := rand.Intn(len(servers))
		addr := servers[idx]
		// Redirect service
		fn = strings.ToLower(fn)
		urlAddr := fmt.Sprintf("http://%s/%s", addr, fn)
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

	ioginx.NewIOServer(svc).Run(addr, config.Name)
}
