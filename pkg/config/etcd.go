package ioconfig

import (
	"log"
	"sync"
	"time"
)

var (
	etcdConf *EtcdConf
	etcdOnce sync.Once
)

type EtcdConf struct {
	Endpoints   []string      `mapstructure:"endpoints"`
	DialTimeout time.Duration `mapstructure:"timeout"` // 连接重试间隔
}

func GetEtcdConf() *EtcdConf {
	etcdOnce.Do(func() {
		if etcdConf == nil {
			if err := V.UnmarshalKey("etcd", &etcdConf); etcdConf == nil || err != nil {
				log.Panicf("unmarshal conf failed, err: %s \n", err)
			}
			etcdConf.DialTimeout = etcdConf.DialTimeout * time.Second
		}
	})
	return etcdConf
}
