package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"fmt"
)

const (
	DES = iota
	AES
)

// 加密
func Encrypt(src []byte, algo int, key []byte) ([]byte, error) {
	var block cipher.Block
	var err error
	switch algo {
	case AES:
		block, err = aes.NewCipher(key)
	case DES:
		block, err = des.NewCipher(key)
	default:
		return nil, fmt.Errorf("unsurported encrypt algo %d", algo)
	}
	// 统一处理错误
	if err != nil {
		fmt.Printf("加密失败，err:%v\n", err)
		return nil, err
	}
	// 填充补充位
	pksc, err := NewPkcs7(block.BlockSize())
	if err != nil {
		fmt.Printf("加密失败，err:%v\n", err)
		return nil, err
	}
	src, err = pksc.Padding(src)
	if err != nil {
		fmt.Printf("加密失败，err:%v\n", err)
		return nil, err
	}
	encrypting := cipher.NewCBCEncrypter(block, key) // 采用 CBC 分组模式加密
	// 加密
	dest := make([]byte, len(src)) // 密文与明文长度相同
	encrypting.CryptBlocks(dest, src)
	return dest, nil
}

// 解密
func Decrypt(in []byte, algo int, key []byte) ([]byte, error) {
	var block cipher.Block
	var err error
	switch algo {
	case AES:
		block, err = aes.NewCipher(key)
	case DES:
		block, err = des.NewCipher(key)
	default:
		return nil, fmt.Errorf("unsurported encrypt algo %d", algo)
	}
	// 统一处理错误
	if err != nil {
		return nil, err
	}
	// 解密
	encrypting := cipher.NewCBCDecrypter(block, key) //采用 CBC 分组模式解密
	transition := make([]byte, len(in))
	encrypting.CryptBlocks(transition, in)
	// 去除填充
	pksc, err := NewPkcs7(block.BlockSize())
	if err != nil {
		return nil, err
	}
	out, err := pksc.UnPadding(transition)
	if err != nil {
		return nil, err
	}
	return out, nil
}
