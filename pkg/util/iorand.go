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
	randomCode := fmt.Sprintf("%06d", int(randomBytes[0])<<8|int(randomBytes[1])<<16|int(randomBytes[2]))
	return randomCode, nil
}
