package ioetcd

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	ioconfig "ims-server/pkg/config"
	iologger "ims-server/pkg/logger"
)

func NewClient() (*clientv3.Client, error) {
	// TODO：设置超时时间
	config := ioconfig.GetEtcdConf()
	client, err := clientv3.New(
		clientv3.Config{
			Endpoints:   config.Endpoints,
			DialTimeout: config.DialTimeout,
		},
	)
	if err != nil {
		iologger.Fatalf("etcd client connect failed, err: %v", err)
	}
	return client, err
}
