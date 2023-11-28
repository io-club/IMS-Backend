package ioetcd

import (
	etcd "go.etcd.io/etcd/client/v3"
	ioconfig "ims-server/pkg/config"
	iologger "ims-server/pkg/logger"
)

func NewClient() (*etcd.Client, error) {
	// TODO：设置超时时间
	config := ioconfig.GetEtcdConf()
	client, err := etcd.New(
		etcd.Config{
			Endpoints:   config.Endpoints,
			DialTimeout: config.DialTimeout,
		},
	)
	if err != nil {
		iologger.Fatalf("etcd client connect failed, err: %v", err)
	}
	return client, err
}
