package ioginx

import (
	"fmt"
	"ims-server/pkg/util"
	"net/http"
	"reflect"
	"strings"
)

type Route struct {
	Func interface{} // Function implementation

	funcName string
	reqName  string // Request struct name
	respName string // Response struct name

	private bool // Whether it is an internal public function

	Permission util.Set[string]
	Methods    []string // Method types,such as: get, post...
}

func (r *Route) GetFuncName() string {
	return r.funcName
}

func ParseRoute(route *Route) error {
	fn := route.Func

	funcType := reflect.TypeOf(fn)
	if funcType.Kind() != reflect.Func {
		return fmt.Errorf("route.Func is not a func")
	}
	if funcType.NumIn() != 2 || funcType.NumOut() != 2 {
		return fmt.Errorf("the function parameters or return values do not meet the requirements")
	}
	if len(route.Methods) == 0 {
		return fmt.Errorf("HTTP methods for the function are not specified")
	}

	route.funcName = util.GetFunctionName(fn)
	route.reqName = funcType.In(1).Name()
	route.respName = funcType.In(0).Name()
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
			lowerUrl := strings.ToLower(r.funcName)
			if _, ok := m[lowerUrl]; ok {
				return fmt.Errorf("duplicate route registration: %s", r.funcName)
			}
			if r.reqName != r.funcName+"Request" {
				return fmt.Errorf("%s request struct name does not meet the requirements, it should be: %s", r.reqName, r.funcName+"Request")
			}
			if r.respName != r.funcName+"Response" {
				return fmt.Errorf("%s response struct name does not meet the requirements, it should be: %s", r.respName, r.funcName+"Response")
			}

			m[method][lowerUrl] = struct{}{}
		}
	}
	return nil
}
