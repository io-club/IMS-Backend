package ioconfig

import (
	"log"
	"sync"
)

var (
	minioConf *MinioConf
	minioOnce sync.Once
)

type MinioConf struct {
	AccessKey     string
	SecretKey     string
	Endpoint      string
	UseSLL        bool
	DefaultBucket string
}

func GetMinioConf() *MinioConf {
	minioOnce.Do(func() {
		if minioConf == nil {
			if err := V.UnmarshalKey("minio", &minioConf); minioConf == nil || err != nil {
				log.Panicf("unmarshal conf failed, err: %s \n", err)
			}
		}
	})
	return minioConf
}
