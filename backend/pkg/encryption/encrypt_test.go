package encryption

import (
	"fmt"
	"testing"
)

func TestEncryptAndDecrypt(t *testing.T) {
	key := []byte("hejunrutianjinpe")
	fmt.Printf("len: %d\n", len(key))
	s := "admin123"
	enc, err := Encrypt([]byte(s), AES, key)
	if err != nil {
		fmt.Printf("加密失败，err:%v\n", err)
		return
	}
	fmt.Printf("加密获得 encrypt: %s\n", string(enc))
	src, err := Decrypt(enc, AES, key)
	if err != nil {
		fmt.Printf("解密失败，err:%v\n", err)
		return
	}
	fmt.Printf("解密获得 src: %s\n", string(src))
}
