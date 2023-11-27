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

func StructToMap(obj interface{}) map[string]interface{} {
	values := reflect.ValueOf(obj)

	if values.Kind() == reflect.Ptr {
		values = values.Elem()
	}
	if values.Kind() != reflect.Struct {
		return nil
	}

	m := make(map[string]interface{})
	typ := values.Type()
	for i := 0; i < values.NumField(); i++ {
		field := values.Field(i)
		name := strings.ToLower(typ.Field(i).Name)
		m[name] = field.Interface()
	}

	return m
}
