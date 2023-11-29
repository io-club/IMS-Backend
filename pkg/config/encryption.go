package ioconfig

import (
	"log"
	"sync"
)

var (
	encryptionConf   = &EncryptionConf{}
	encryptionConfig *EncryptionConfig
	encryptionOnce   sync.Once
)

type EncryptionConfig struct {
	AesKey string `mapstructure:"aesKey"`
	DesKey string `mapstructure:"desKey"`
}

type EncryptionConf struct {
	AesKey []byte
	DesKey []byte
}

func GetEncryptionConf() *EncryptionConf {
	encryptionOnce.Do(func() {
		if encryptionConfig == nil {
			if err := V.UnmarshalKey("encryption", &encryptionConfig); encryptionConfig == nil || err != nil {
				log.Panicf("unmarshal conf failed, err: %v\n", err)
			}
			encryptionConf.AesKey = []byte(encryptionConfig.AesKey)
			encryptionConf.DesKey = []byte(encryptionConfig.DesKey)
		}
	})
	return encryptionConf
}
