package ioetcd

import (
	etcd "go.etcd.io/etcd/client/v3"
	ioconfig "ims-server/pkg/config"
	iologger "ims-server/pkg/logger"
	"time"
)

func NewClient() (*etcd.Client, error) {
	var client *etcd.Client
	var err error
	connectDoneCh := make(chan struct{})
	go func() {
		config := ioconfig.GetEtcdConf()
		client, err = etcd.New(
			etcd.Config{
				Endpoints:   config.Endpoints,
				DialTimeout: config.DialTimeout,
			},
		)
		if err != nil {
			iologger.Fatalf("etcd client connect failed, err: %v", err)
		}
		connectDoneCh <- struct{}{}
	}()
	select {
	case <-time.After(5 * time.Second):
		// 设置 5 秒的超时限制
		iologger.Fatalf("etcd client connect timeout")
	case <-connectDoneCh:
	}
	return client, err
}
