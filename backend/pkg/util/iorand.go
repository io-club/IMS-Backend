package util

import (
	"crypto/rand"
	"fmt"
)

func GetRandCode() (string, error) {
	// 随机获取 3 个字节
	randomBytes := make([]byte, 3)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	// 生成个 6 位数字
	randomCode := fmt.Sprintf("%06d", int(randomBytes[0])<<6|int(randomBytes[1])<<5|int(randomBytes[2]<<3))
	return randomCode, nil
}
