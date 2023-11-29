package encryption

import (
	"bytes"
	"errors"
)

type pkcs interface {
	Padding(src []byte, blockSize int) []byte
	UnPadding(src []byte, blockSize int) ([]byte, error)
}

// block 可以是 1~255
type pkcs7 struct {
	pkcs
	blockSize int
}

func NewPkcs7(blockSize int) (*pkcs7, error) {
	if blockSize < 1 || blockSize > 255 {
		return nil, errors.New("pkcs7 不支持该 blockSize")
	}
	return &pkcs7{blockSize: blockSize}, nil
}

func (p *pkcs7) Padding(src []byte) ([]byte, error) {
	srcLen := len(src)
	padLen := p.blockSize - (srcLen % p.blockSize)
	// 生成填充字符串
	padText := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(src, padText...), nil
}

func (p *pkcs7) UnPadding(src []byte) ([]byte, error) {
	srcLen := len(src)
	paddingLen := int(src[srcLen-1])
	if srcLen < paddingLen || paddingLen > p.blockSize {
		return nil, errors.New("blockSize 大小错误")
	}
	return src[:srcLen-paddingLen], nil
}

// pkcs5 是 pkcs7 的子集，与 pkcs7 的区别只在于 block 必须为 8
type pkcs5 struct {
	pkcs
	blockSize int
}

func NewPkcs5() *pkcs5 {
	return &pkcs5{blockSize: 8}
}

func (p *pkcs5) Padding(src []byte) []byte {
	blockSize := 8
	srcLen := len(src)
	paddingLen := blockSize - srcLen%blockSize
	padText := bytes.Repeat([]byte{byte(paddingLen)}, paddingLen)
	return append(src, padText...)
}

func (p *pkcs5) UnPadding(src []byte) ([]byte, error) {
	blockSize := 8
	srcLen := len(src)
	paddingLen := int(src[srcLen-1])
	if srcLen < paddingLen || blockSize < paddingLen {
		return nil, errors.New("传入的数据不是 pkcs5 标准")
	}
	return src[:srcLen-paddingLen], nil
}
