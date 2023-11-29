package main

import (
	userroute "ims-server/internal/user/route"
	ioconfig "ims-server/pkg/config"
	ioginx "ims-server/pkg/ginx"
	"strings"
)

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

var NmsRoutesMap = map[string][]ioginx.Route{
	ioconfig.GetServiceConf().User.Name: userroute.Routes,
}
