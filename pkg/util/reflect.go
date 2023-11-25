package util

import (
	"reflect"
	"runtime"
	"strings"
)

func GetFunctionName(f interface{}) string {
	// Get the pointer value of the function
	funcPtr := reflect.ValueOf(f).Pointer()

	// Get the function object
	funcObj := runtime.FuncForPC(funcPtr)

	// Get the function name
	funcName := funcObj.Name()

	// Use string manipulation to extract the last part of the name
	parts := strings.Split(funcName, ".")
	funcLastName := parts[len(parts)-1]

	return funcLastName
}
