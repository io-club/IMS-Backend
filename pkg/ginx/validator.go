package ioginx

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

// 自定义密码验证函数
func CheckPassword(fl validator.FieldLevel) bool {
	password, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	// 正则表达式要求密码至少包含一个字母和一个数字
	pattern := `^[a-zA-Z0-9]+$`
	ok, err := regexp.MatchString(pattern, password)
	if err != nil {
		return false
	}
	return ok
}
