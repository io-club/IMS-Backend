package ioginx

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

// 自定义密码验证函数
func customPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// 使用正则表达式验证密码必须包含字母和数字
	// 这个正则表达式要求密码至少包含一个字母和一个数字
	pattern := `^(?=.*[a-zA-Z])(?=.*\d).+$`
	ok, err := regexp.MatchString(pattern, password)
	if err != nil {
		return false
	}
	return ok
}
