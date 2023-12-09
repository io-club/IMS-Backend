package ioginx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	ioconsts "ims-server/pkg/consts"
	"ims-server/pkg/util"
	"net/http"
	"reflect"
	"strings"
)

type Route struct {
	Func func(c *gin.Context) // Function implementation

	FuncName string

	Private bool // Whether it is an internal public function

	Permission util.Set[ioconsts.UserType]
	Methods    []string // Method types,such as: get, post...
}

func GetPublicRouteMap(routes []Route) map[string]Route {
	publicRouteMap := map[string]Route{}
	for _, route := range routes {
		if route.Private {
			continue
		}

		if err := ParseRoute(&route); err != nil {
			panic(err)
		}

		publicRouteMap[strings.ToLower(route.FuncName)] = route
	}
	return publicRouteMap
}

func ParseRoute(route *Route) error {
	fn := route.Func

	funcType := reflect.TypeOf(fn)
	if funcType.Kind() != reflect.Func {
		return fmt.Errorf("route.Func is not a func")
	}
	if len(route.Methods) == 0 {
		return fmt.Errorf("HTTP methods for the function are not specified")
	}

	return nil
}

func CheckRoutes(routes []Route) error {
	m := map[string]map[string]struct{}{
		http.MethodGet:  {},
		http.MethodPost: {},
	}
	for _, r := range routes {
		for _, method := range r.Methods {
			// Check for duplication
			lowerUrl := strings.ToLower(r.FuncName)
			if _, ok := m[lowerUrl]; ok {
				return fmt.Errorf("duplicate route registration: %s", r.FuncName)
			}
			m[method][lowerUrl] = struct{}{}
		}
	}
	return nil
}
