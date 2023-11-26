package util

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/iancoleman/strcase"
)

// TODO：完全重构该包
type structToMap struct {
	filterKeySet  Set[string]
	keyConverters []func(string) string
	valueFilters  []func(v reflect.Value) bool
}

// 结构体转 Map
//
// @param: data: 需要转换的结构体
// @param: keyConverter: key 转换函数，如果为 nil，则不转换
// @param: ignoreStrategy: 忽略策略，如果为 nil，则不忽略
// @return: map[string]interface{} 转换后的 map
func NewStructToMap() structToMap {
	return structToMap{
		filterKeySet:  NewSet[string](),
		valueFilters:  []func(v reflect.Value) bool{},
		keyConverters: []func(string) string{},
	}
}

func (s structToMap) WithKeyExclude(keys ...string) structToMap {
	s.filterKeySet.Add(keys...)
	return s
}

func (s structToMap) WithKeyConverter(keyConverters ...func(string string) string) structToMap {
	s.keyConverters = append(s.keyConverters, keyConverters...)
	return s
}

func (s structToMap) WithValueFilter(filters ...func(v reflect.Value) bool) structToMap {
	s.valueFilters = append(s.valueFilters, filters...)
	return s
}

func (s structToMap) Run(data interface{}) map[string]interface{} {
	structValue := reflect.ValueOf(data)
	structType := structValue.Type()

	// 如果传进来的结构体是指针
	// TODO: 如果传进来的是 **SomeStruct
	for structType.Kind() == reflect.Ptr {
		structValue = structValue.Elem()
		structType = structValue.Type()
	}

	result := make(map[string]interface{})
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)
		fieldName := field.Name

		// 过滤键
		if s.filterKeySet.Contains(fieldName) {
			goto nextRound
		}

		//过滤值
		for _, filter := range s.valueFilters {
			if filter(fieldValue) {
				goto nextRound
			}
		}

		// TODO: 如果是指针类型，需要取出指针指向的值

		// 转换键的格式
		for _, keyConverter := range s.keyConverters {
			fieldName = keyConverter(fieldName)
		}

		// 添加到返回值
		result[fieldName] = fieldValue.Interface()

	nextRound:
	}

	return result
}

func FilterEmptyStringStrategy(v reflect.Value) bool {
	if v.Kind() == reflect.String {
		return v.String() == ""
	}
	return false
}
func FilterZeroValueStrategy(v reflect.Value) bool {
	return v.IsZero()
}

// 将请求转换为蛇形的 map，忽略空字符串
func RequestToSnakeMapWithIgnoreEmptyString(data interface{}) map[string]interface{} {
	return NewStructToMap().
		WithKeyConverter(strcase.ToSnake).
		WithValueFilter(FilterEmptyStringStrategy).
		Run(data)
}

func RequestToSnakeMapWithIgnoreZeroValueAndModelKey(data interface{}) map[string]interface{} {
	keys := []string{}
	keys = append(keys, "Model")
	return NewStructToMap().
		WithKeyConverter(strcase.ToSnake).
		WithKeyExclude(keys...).
		WithValueFilter(FilterZeroValueStrategy).
		Run(data)
}

func RequestToSnakeMapWithIgnoreModel(data interface{}) map[string]interface{} {
	keys := []string{"Model"}
	return NewStructToMap().
		WithKeyConverter(strcase.ToSnake).
		WithKeyExclude(keys...).
		Run(data)
}

func RequestToSnakeMapWithIgnoreZeroValueAndIDKey(data interface{}) map[string]interface{} {
	keys := []string{"id"}
	return NewStructToMap().
		WithKeyConverter(strcase.ToSnake).
		WithKeyExclude(keys...).
		WithValueFilter(FilterZeroValueStrategy).
		Run(data)
}

func StructFilterUnCheck[D any](srcUnknown interface{}) D {
	res, _ := StructFilter[D](srcUnknown)
	return res
}

func StructFilter[D any](srcUnknown interface{}) (D, error) {
	srcV := reflect.ValueOf(srcUnknown)
	srcT := srcV.Type()

	// 如果传进来的结构体是指针，则取出指针指向的值
	for srcT.Kind() == reflect.Ptr {
		srcV = srcV.Elem()
		srcT = srcV.Type()
	}

	var dest D

	// 检查源结构体类型
	if srcT.Kind() != reflect.Struct {
		return dest, errors.New("can't filter non struct type")
	}
	// 检查目标结构体类型
	destT := reflect.TypeOf(dest)
	if destT.Kind() != reflect.Struct {
		return dest, errors.New("can't filter to non struct type")
	}

	// 初始化目标结构体
	destV := reflect.New(destT)
	fmt.Printf("[%+v]\n", destV)
	// 遍历目标结构体的字段
	for i := 0; i < destT.NumField(); i++ {
		destField := destT.Field(i)
		// 检查源结构体是否有该字段

		srcField, ok := srcT.FieldByName(destField.Name)
		if !ok {
			continue
		}
		// 检查字段类型一致
		if srcField.Type != destField.Type {
			continue
		}
		// 检查字段是否可写
		if !destV.Elem().Field(i).CanSet() {
			continue
		}
		// 赋值
		destV.Elem().Field(i).Set(srcV.FieldByName(destField.Name))
	}
	res, err := destV.Interface().(*D)
	if !err {
		return dest, errors.New("can't convert to dest type")
	}
	return *res, nil
}

func ExtractColumn[D any](srcUnknownList interface{}, fieldName string) ([]D, error) {
	srcV := reflect.ValueOf(srcUnknownList)
	srcT := srcV.Type()
	// 如果传进来的结构体是指针，则取出指针指向的值
	for srcT.Kind() == reflect.Ptr {
		srcV = srcV.Elem()
		srcT = srcV.Type()
	}
	if srcT.Kind() != reflect.Slice {
		return nil, errors.New("can't extract from non array value")
	}

	res := make([]D, 0, srcV.Len())
	for i := 0; i < srcV.Len(); i++ {
		item := srcV.Index(i)
		if item.Kind() != reflect.Struct {
			return nil, errors.New("item not a struct")
		}
		v := reflect.Indirect(item)
		if v.IsZero() {
			return nil, errors.New("struct should not be zero")
		}
		fieldValue := v.FieldByName(fieldName)
		if fieldValue.IsZero() {
			return nil, errors.New("struct field should not be zero")
		}
		target, ok := fieldValue.Interface().(D)
		if !ok {
			return nil, errors.New("struct field not valid: " + fmt.Sprintf("%v", reflect.TypeOf(fieldValue.Interface()).Kind()))
		}
		res = append(res, target)
	}
	return res, nil
}

func ExtractColumnAsMapKey[K comparable, V any](srcUnknownList []V, fieldName string) (map[K]V, error) {
	res := make(map[K]V, len(srcUnknownList))
	for i := 0; i < len(srcUnknownList); i++ {
		item := srcUnknownList[i]
		itemT := reflect.TypeOf(item)
		if itemT.Kind() != reflect.Struct {
			return nil, errors.New("item not a struct")
		}
		itemV := reflect.ValueOf(item)
		v := reflect.Indirect(itemV)
		if v.IsZero() {
			return nil, errors.New("struct should not be zero")
		}
		fieldValue := v.FieldByName(fieldName)
		if fieldValue.IsZero() {
			return nil, errors.New("struct field should not be zero")
		}
		target, ok := fieldValue.Interface().(K)
		if !ok {
			return nil, errors.New("struct field not valid: " + fmt.Sprintf("%v", reflect.TypeOf(fieldValue.Interface()).Kind()))
		}
		res[target] = item
	}
	return res, nil
}
